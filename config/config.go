// Package config implements tag's config file functionality.
//
// Tag can operate perfectly well without a config file, but it's
// used to specify hooks, custom commit and tag messages and to
// configure search and replace on version bumping for e.g. project metadata
package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/pelletier/go-toml/v2"
)

// initContents is the contents of the initial config file created by `tag init`
//
//go:embed init.toml
var initContents string

// Filename is the canonical file name for tag's config file.
const Filename = ".tag.toml"

// ErrNoConfigFile is the signal that no config file exists.
var ErrNoConfigFile = errors.New("config file not found")

const filePermissions = 0o644

// Config represents tags configuration settings.
type Config struct { //nolint: recvcheck // In this case it makes sense
	Version string `toml:"version"`
	Git     Git    `toml:"git,omitempty"`
	Hooks   Hooks  `toml:"hooks,omitempty"`
	Files   []File `toml:"file,omitempty"`
}

// Git represents the git config in tag's config file.
type Git struct {
	DefaultBranch   string `toml:"default-branch,omitempty"`
	MessageTemplate string `toml:"message-template,omitempty"`
	TagTemplate     string `toml:"tag-template,omitempty"`
}

// Hooks encodes the optional hooks specified in tag's config file.
type Hooks struct {
	PreReplace string `toml:"pre-replace,omitempty"`
	PreCommit  string `toml:"pre-commit,omitempty"`
	PreTag     string `toml:"pre-tag,omitempty"`
	PrePush    string `toml:"pre-push,omitempty"`
}

// File represents a single file tag should perform search and replace on.
type File struct {
	Path   string `toml:"path,omitempty"`
	Search string `toml:"search,omitempty"`

	Replace string `toml:"-"` // Not part of the config, inferred from `Search`
}

// Load reads Config from a file.
func Load(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrNoConfigFile
		}
		return Config{}, fmt.Errorf("could not read %s: %w", path, err)
	}

	if len(bytes.TrimSpace(raw)) == 0 {
		return Config{}, fmt.Errorf("config file %s is empty", path)
	}

	cfg := Config{
		// Git has a default
		Git: Git{
			DefaultBranch:   "main",
			MessageTemplate: "Bump version {{.Current}} -> {{.Next}}",
			TagTemplate:     "v{{.Next}}",
		},
	}
	if err := toml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, fmt.Errorf("toml deserialize error: %w", err)
	}

	return cfg, nil
}

// Save saves the config to disk.
func (c Config) Save(path string) error {
	raw, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("toml serialise error: %w", err)
	}
	if err := os.WriteFile(path, raw, filePermissions); err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}
	return nil
}

// Render replaces the special values {{.Current}} and {{.Next}} in the
// search and replace templates as well as the commit and tag messages.
func (c *Config) Render(current, next string) error {
	searchTemplate := template.New("search")
	replaceTemplate := template.New("replace")
	tagTemplate := template.New("tag")
	commitTemplate := template.New("commit")

	vars := map[string]string{"Current": current, "Next": next}

	// Render the tag and commit templates
	tagParsed, err := tagTemplate.Parse(c.Git.TagTemplate)
	if err != nil {
		return fmt.Errorf("could not parse tag-template: %w", err)
	}
	commitParsed, err := commitTemplate.Parse(c.Git.MessageTemplate)
	if err != nil {
		return fmt.Errorf("could not parse message-template: %w", err)
	}

	tagOut := &bytes.Buffer{}
	commitOut := &bytes.Buffer{}

	if err := tagParsed.Execute(tagOut, vars); err != nil {
		return fmt.Errorf("could not execute tag-template: %w", err)
	}

	if err := commitParsed.Execute(commitOut, vars); err != nil {
		return fmt.Errorf("could not execute message-template: %w", err)
	}

	// Overwrite the originals with the now-rendered text
	c.Git.TagTemplate = tagOut.String()
	c.Git.MessageTemplate = commitOut.String()

	rendered := make([]File, 0, len(c.Files))

	// Now for the files
	for _, file := range c.Files {
		searchParsed, err := searchTemplate.Parse(file.Search)
		if err != nil {
			return fmt.Errorf("could not parse file.search for file %s: %w", file.Path, err)
		}

		replaceParsed, err := replaceTemplate.Parse(strings.ReplaceAll(file.Search, "{{.Current}}", "{{.Next}}"))
		if err != nil {
			return fmt.Errorf("could not infer file.replace for file %s: %w", file.Path, err)
		}

		searchOut := &bytes.Buffer{}
		replaceOut := &bytes.Buffer{}
		if err := searchParsed.Execute(searchOut, vars); err != nil {
			return fmt.Errorf("could not execute file.search for file %s: %w", file.Path, err)
		}

		if err := replaceParsed.Execute(replaceOut, vars); err != nil {
			return fmt.Errorf("could not execute file.replace for file %s: %w", file.Path, err)
		}

		file.Search = searchOut.String()
		file.Replace = replaceOut.String()
		rendered = append(rendered, file)
	}

	c.Files = rendered
	return nil
}

// Init returns a toml encoded string of the initial tag config.
func Init() string {
	return initContents
}
