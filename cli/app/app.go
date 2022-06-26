// Package app implements the CLI functionality, the CLI defers
// execution to the exported methods in this package
package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/tag/config"
	"github.com/FollowTheProcess/tag/git"
	"github.com/FollowTheProcess/tag/replacer"
	"github.com/FollowTheProcess/tag/version"
)

const (
	major  = "major"
	minor  = "minor"
	patch  = "patch"
	master = "master"
	main   = "main"
	trunk  = "trunk"
)

var errAbort = errors.New("Aborted")

// App represents the tag program.
type App struct {
	out     io.Writer      // Where to write output to
	printer *msg.Printer   // The app's printer
	config  *config.Config // The tag config
	replace bool           // Whether or not we want to do search and replace
}

// New creates and returns a new app writing to 'out'
// and using the config file at 'path', if the config file
// does not exist, app.Replace will be false.
func New(out io.Writer, path string) *App {
	printer := msg.Default()
	printer.Out = out

	app := &App{out: out, printer: printer}

	cwd, err := os.Getwd()
	if err != nil {
		// Don't really like panicking but if we can't
		// even get the cwd it's probably the only thing to do
		// as New can't return an error
		panic(err)
	}

	path = filepath.Join(cwd, path)

	cfg, err := config.Load(path)
	if err != nil {
		app.replace = false
	} else {
		app.printer.Infof("Config file %s found and loaded", path)
		app.config = cfg
		app.replace = true
	}

	return app
}

// Patch is the tag patch subcommand.
func (a *App) Patch(force, push bool, message string) error {
	return a.bump(patch, message, force, push)
}

// Minor is the tag minor subcommand.
func (a *App) Minor(force, push bool, message string) error {
	return a.bump(minor, message, force, push)
}

// Major is the tag major subcommand.
func (a *App) Major(force, push bool, message string) error {
	return a.bump(major, message, force, push)
}

// List is what happens when tag is invoked with no subcommands.
func (a *App) List() error {
	list, err := git.ListTags()
	if err != nil {
		return err
	}
	fmt.Fprint(a.out, list)
	return nil
}

// Latest is the tag latest subcommand.
func (a *App) Latest() error {
	latest, err := git.LatestTag()
	if err != nil {
		return err
	}
	fmt.Fprintln(a.out, latest)
	return nil
}

// EnsureGitRepo returns an error if tag is invoked outside
// of a git repo (other than --help or --version).
func (a *App) EnsureGitRepo() error {
	if !git.IsRepo() {
		return errors.New("Not inside a git repo")
	}
	return nil
}

// ensureBumpableState returns an error if the working tree
// is dirty or we're not on a branch called master, main, or trunk.
func (a *App) ensureBumpableState() error {
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

	if (branch != main) && (branch != master) && (branch != trunk) {
		return fmt.Errorf("Must be on branch {main | master | trunk} to bump a version, currently on %s", branch)
	}
	return nil
}

// bump is the generic bump function, that holds a lot of shared setup
// and dispatches to the correct type.
func (a *App) bump(typ, message string, force, push bool) error {
	if err := a.ensureBumpableState(); err != nil {
		return err
	}
	latest, err := git.LatestTag()
	if err != nil {
		return err
	}
	current, err := version.Parse(latest)
	if err != nil {
		return err
	}

	var next version.Version

	switch typ {
	case major:
		next = version.BumpMajor(current)
	case minor:
		next = version.BumpMinor(current)
	case patch:
		next = version.BumpPatch(current)
	default:
		panic("unreachable")
	}

	if !force {
		confirm := &survey.Confirm{
			Message: fmt.Sprintf("This will bump %q to %q. Are you sure?", current, next),
			Default: false,
		}
		err = survey.AskOne(confirm, &force)
		if err != nil {
			return err
		}
	}

	// Now if force is false, the user said no -> abort
	if !force {
		return errAbort
	}

	if err := a.doBump(current, next, message, push); err != nil {
		return err
	}

	a.printer.Good("Done")
	return nil
}

// doBump is a helper that actually does the bumping (including any replacing).
func (a *App) doBump(current, next version.Version, message string, push bool) error {
	if a.replace {
		if err := a.config.Render(current.String(), next.String()); err != nil {
			return err
		}
		for _, file := range a.config.Files {
			// TODO: If there are multiple entries for the same file, this will open
			// and close it multiple times which is not ideal
			err := replacer.Replace(file.Path, file.Search, file.Replace)
			if err != nil {
				return err
			}
		}
		if err := git.Add(); err != nil {
			return err
		}
		stdout, err := git.Commit(fmt.Sprintf("Bump version %s -> %s", current.String(), next.String()))
		if err != nil {
			return errors.New(stdout)
		}
	}

	stdout, err := git.CreateTag(next.Tag(), message)
	if err != nil {
		return errors.New(stdout)
	}

	// If push, push the tag
	if push {
		if a.replace {
			a.printer.Info("Pushing bump commit")
			stdout, err = git.Push()
			if err != nil {
				return errors.New(stdout)
			}
		}
		a.printer.Infof("Pushing tag %s", next.Tag())
		stdout, err = git.PushTag(next.Tag())
		if err != nil {
			return errors.New(stdout)
		}
	}
	return nil
}
