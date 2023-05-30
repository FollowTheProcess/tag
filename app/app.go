// Package app implements the functionality of tag, the CLI calls
// exported members of this package.
package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/semver"
	"github.com/FollowTheProcess/tag/config"
	"github.com/FollowTheProcess/tag/git"
	"github.com/FollowTheProcess/tag/hooks"
)

// ErrAborted is returned whenever an action is aborted by the user.
var ErrAborted = errors.New("Aborted")

// App represents the tag program.
type App struct {
	Stdout      io.Writer
	Stderr      io.Writer
	Cfg         config.Config
	replaceMode bool
}

// bumpType is an enum of recognised bump types.
type bumpType int

const (
	major bumpType = iota
	minor
	patch
)

// New constructs and returns a new App.
func New(cwd string, stdout, stderr io.Writer) (App, error) {
	path := filepath.Join(cwd, config.Filename)
	replaceMode := true
	cfg, err := config.Load(path)
	if err != nil {
		if errors.Is(err, config.ErrNoConfigFile) {
			replaceMode = false
		} else {
			return App{}, err
		}
	}

	app := App{
		Stdout:      stdout,
		Stderr:      stderr,
		Cfg:         cfg,
		replaceMode: replaceMode,
	}

	return app, nil
}

// List handles the list subcommand.
func (a App) List(limit int) error {
	if err := a.ensureRepo(); err != nil {
		return err
	}
	if limit <= 0 {
		return fmt.Errorf("--limit must be a positive integer")
	}
	tags, limitHit, err := git.ListTags(limit)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Stdout, strings.TrimSpace(tags))
	if limitHit {
		fmt.Fprintln(a.Stdout)
		msg.Fwarn(a.Stdout, "Truncated, pass --limit to see more")
	}

	return nil
}

// Latest handles the latest subcommand.
func (a App) Latest() error {
	if err := a.ensureRepo(); err != nil {
		return err
	}
	tag, err := git.LatestTag()
	if err != nil {
		return err
	}
	fmt.Fprintln(a.Stdout, tag)
	return nil
}

// Init handles the init subcommand.
func (a App) Init(cwd string, force bool) error {
	path := filepath.Join(cwd, config.Filename)
	configFileExists, err := exists(path)
	if err != nil {
		return err
	}

	cfg := config.Config{
		Version: "0.1.0",
		Git: config.Git{
			DefaultBranch:   "main",
			MessageTemplate: "Bump version {{.Current}} -> {{.Next}}",
			TagTemplate:     "v{{.Next}}",
		},
		Hooks: config.Hooks{
			PreReplace: "echo 'I run before doing anything'",
			PreCommit:  "echo 'I run after replacing but before committing changes'",
			PreTag:     "echo 'I run after committing changes but before tagging'",
			PrePush:    "echo 'I run after tagging, but before pushing'",
		},
		Files: []config.File{
			{
				Path:   "pyproject.toml",
				Search: `version = "{{.Current}}"`,
			},
			{
				Path:   "README.md",
				Search: "My project, version {{.Current}}",
			},
		},
	}

	if !configFileExists {
		// No config file, just go ahead and make one
		if err := cfg.Save(path); err != nil {
			return err
		}
		msg.Fsuccess(a.Stdout, "Config file written to %s", path)
		return nil
	}

	// Config file does exist, let's ask for overwrite and check force
	if !force {
		confirm := &survey.Confirm{
			Message: fmt.Sprintf("Config file %s already exists. Overwrite?", path),
			Default: false,
		}
		err := survey.AskOne(confirm, &force)
		if err != nil {
			return err
		}
	}

	// Now if force is still false, user said no -> abort
	if !force {
		return ErrAborted
	}

	// User has either confirmed or passed --force
	if err := cfg.Save(path); err != nil {
		return err
	}
	msg.Fsuccess(a.Stdout, "Config file written to %s", path)
	return nil
}

// Major handles the major subcommand.
func (a App) Major(push, force, dryRun bool) error {
	return a.bump(major, push, force, dryRun)
}

// Minor handles the minor subcommand.
func (a App) Minor(push, force, dryRun bool) error {
	return a.bump(minor, push, force, dryRun)
}

// Patch handles the minor subcommand.
func (a App) Patch(push, force, dryRun bool) error {
	return a.bump(patch, push, force, dryRun)
}

// replaceAll is a helper that performs and reports on file replacement
// as part of bumping.
func (a App) replaceAll(current, next semver.Version, dryRun bool) error {
	if err := a.Cfg.Render(current.String(), next.String()); err != nil {
		return err
	}

	if err := a.replace(dryRun); err != nil {
		return err
	}

	// Also replace the Version in the config file
	a.Cfg.Version = next.String()
	if !dryRun {
		if err := a.Cfg.Save(config.Filename); err != nil {
			return err
		}
	}

	dirty, err := git.IsDirty()
	if err != nil {
		return err
	}

	if err = a.runHook(hooks.StagePreCommit, dryRun); err != nil {
		return err
	}

	// Only any point in committing if something has changed
	if dirty {
		if dryRun {
			msg.Finfo(a.Stdout, "(Dry Run) Would commit changes")
			return nil
		}
		msg.Finfo(a.Stdout, "Committing changes")
		if err = git.Add(); err != nil {
			return err
		}

		commitOut, err := git.Commit(a.Cfg.Git.MessageTemplate)
		if err != nil {
			return errors.New(commitOut)
		}
	}
	return nil
}

// replace is a helper that performs file replacement.
func (a App) replace(dryRun bool) error {
	for _, file := range a.Cfg.Files {
		var contents []byte
		contents, err := os.ReadFile(file.Path)
		if err != nil {
			return err
		}

		if !bytes.Contains(contents, []byte(file.Search)) {
			return fmt.Errorf("Could not find %q in %s", file.Search, file.Path)
		}

		if dryRun {
			msg.Finfo(a.Stdout, "(Dry Run) Would replace %s with %s in %s", file.Search, file.Replace, file.Path)
		} else {
			msg.Finfo(a.Stdout, "Replacing contents in %s", file.Path)
			newContent := bytes.ReplaceAll(contents, []byte(file.Search), []byte(file.Replace))

			if err = os.WriteFile(file.Path, newContent, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

// getBumpVersions is a helper that gets .Current and .Next from context.
func (a App) getBumpVersions(typ bumpType) (current, next semver.Version, err error) {
	if a.replaceMode {
		// If the config file is present, use the version specified in there
		current, err = semver.Parse(a.Cfg.Version)
		if err != nil {
			return semver.Version{}, semver.Version{}, err
		}
	} else {
		// Otherwise start at the latest semver tag present
		latest, err := git.LatestTag()
		if err != nil {
			if errors.Is(err, git.ErrNoTagsFound) {
				current = semver.Version{} // No tags, no default version, start at v0.0.0
			} else {
				return semver.Version{}, semver.Version{}, err
			}
		} else {
			current, err = semver.Parse(latest)
			if err != nil {
				return semver.Version{}, semver.Version{}, err
			}
		}
	}

	switch typ {
	case major:
		next = semver.BumpMajor(current)
	case minor:
		next = semver.BumpMinor(current)
	case patch:
		next = semver.BumpPatch(current)
	default:
		return semver.Version{}, semver.Version{}, fmt.Errorf("Unrecognised bump type: %v", typ)
	}

	return current, next, nil
}

// bump is a helper that performs logic common to all bump methods.
func (a App) bump(typ bumpType, push, force, dryRun bool) error {
	if err := a.ensureRepo(); err != nil {
		return err
	}
	if err := a.ensureBumpable(); err != nil {
		return err
	}

	current, next, err := a.getBumpVersions(typ)
	if err != nil {
		return err
	}

	if !force {
		confirm := &survey.Confirm{
			Message: fmt.Sprintf("This will bump %q to %q. Are you sure?", current, next),
			Default: false,
		}
		err := survey.AskOne(confirm, &force)
		if err != nil {
			return err
		}
	}

	// Now if force is false, the user said no -> abort
	if !force {
		return ErrAborted
	}

	if err := a.runHook(hooks.StagePreReplace, dryRun); err != nil {
		return err
	}

	if a.replaceMode {
		if err := a.replaceAll(current, next, dryRun); err != nil {
			return err
		}
	}

	if err := a.runHook(hooks.StagePreTag, dryRun); err != nil {
		return err
	}

	if dryRun {
		msg.Finfo(a.Stdout, "(Dry Run) Would issue new tag %s", next.Tag())
	} else {
		msg.Finfo(a.Stdout, "Issuing new tag %s", next.Tag())
		stdout, err := git.CreateTag(next.Tag(), a.Cfg.Git.TagTemplate)
		if err != nil {
			return errors.New(stdout)
		}
	}

	// If --push, push the tag and commit
	if push {
		if err := a.runHook(hooks.StagePrePush, dryRun); err != nil {
			return err
		}
		if dryRun {
			msg.Finfo(a.Stdout, "(Dry Run) Would push tag %s", next.Tag())
			return nil
		}
		msg.Finfo(a.Stdout, "Pushing tag %s", next.Tag())
		stdout, err := git.Push()
		if err != nil {
			return errors.New(stdout)
		}
	}
	return nil
}

// runHook is a helper that runs a particular hook stage (if it is defined)
// and understands --dry-run.
func (a App) runHook(stage hooks.HookStage, dryRun bool) error {
	var hookCmd string
	switch stage {
	case hooks.StagePreReplace:
		hookCmd = a.Cfg.Hooks.PreReplace
	case hooks.StagePreCommit:
		hookCmd = a.Cfg.Hooks.PreCommit
	case hooks.StagePreTag:
		hookCmd = a.Cfg.Hooks.PreTag
	case hooks.StagePrePush:
		hookCmd = a.Cfg.Hooks.PrePush
	default:
		return fmt.Errorf("Unhandled hook type: %s", stage)
	}

	if hookCmd == "" {
		// No op if the hook is not defined
		return nil
	}

	if dryRun {
		msg.Finfo(a.Stdout, "(Dry Run) Would run hook %s: %s", stage, hookCmd)
		return nil
	}
	return hooks.Run(stage, hookCmd, a.Stdout, a.Stderr)
}

// ensureRepo is a helper that will error if the current directory is not
// a git repo.
func (a App) ensureRepo() error {
	if !git.IsRepo() {
		return errors.New("Not a git repo")
	}
	return nil
}

// ensureBumpable is a helper that will error if the current git state is not
// "bumpable", that is we're on the default branch, and the working tree is clean.
func (a App) ensureBumpable() error {
	dirty, err := git.IsDirty()
	if err != nil {
		return err
	}
	if dirty {
		return errors.New("Working tree is not clean")
	}

	branch, err := git.Branch()
	if err != nil {
		return err
	}

	if a.Cfg.Git.DefaultBranch == "" {
		a.Cfg.Git.DefaultBranch = "main" // Default
	}

	if branch != a.Cfg.Git.DefaultBranch {
		return fmt.Errorf("Not on default branch (%s), currently on: %s", a.Cfg.Git.DefaultBranch, branch)
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
