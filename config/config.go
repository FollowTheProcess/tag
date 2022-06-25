// Package config implements tag's configuration file (.tag.toml)
// tag can operate perfectly well without a config file, but it's
// used to specify files where tag should perform search and replace
// on version bumping.
//
// Example config:
// [tag]
// files = [
//		{ path = "path/to/file.go", search = "version = {{.Current}}", replace = "version = {{.Next}}"},
//		{ path = "README.md", search = "version = {{.Current}}", replace = "version = {{.Next}}"}
// ]
package config

import (
	"bytes"
	"os"
	"text/template"

	"github.com/pelletier/go-toml/v2"
)

const Path = ".tag.toml" // The canonical config file path (relative to cwd).

// File represents a file that tag should perform replacement on.
type File struct {
	Path    string `toml:"path"`    // The path to the file (relative to .tag.toml)
	Search  string `toml:"search"`  // String to search for
	Replace string `toml:"replace"` // String to replace 'Search' with
}

// Tag represents the top level config of the tag program.
type Tag struct {
	Files []File `toml:"files"` // List of files tag should operate on
}

// Config represents tag's configuration.
type Config struct {
	Tag `toml:"tag"` // Top level config
}

// Render replaces the special values {{.Current}} and {{.Next}} in the user's
// configured search and replace directives.
func (c *Config) Render(current, next string) error {
	searchTemplate := template.New("search")
	replaceTemplate := template.New("replace")

	vars := map[string]string{"Current": current, "Next": next}
	rendered := make([]File, 0, len(c.Files))

	for _, file := range c.Files {
		searchParsed, err := searchTemplate.Parse(file.Search)
		if err != nil {
			return err
		}
		replaceParsed, err := replaceTemplate.Parse(file.Replace)
		if err != nil {
			return err
		}

		searchOut := &bytes.Buffer{}
		replaceOut := &bytes.Buffer{}

		if err := searchParsed.Execute(searchOut, vars); err != nil {
			return err
		}

		if err := replaceParsed.Execute(replaceOut, vars); err != nil {
			return err
		}

		// Overwrite the originals with the now rendered text
		file.Search = searchOut.String()
		file.Replace = replaceOut.String()

		rendered = append(rendered, file)
	}

	c.Files = rendered
	return nil
}

// Load reads in tag's config file.
func Load(path string) (*Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := toml.Unmarshal(contents, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
