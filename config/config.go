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
	"os"

	"github.com/pelletier/go-toml/v2"
)

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

// Load reads in tag's config file.
func Load(path string) (Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := toml.Unmarshal(contents, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
