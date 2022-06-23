// Package git implements tags interface with git in order to interact
// with tags and make commits.
package git

import (
	"os/exec"
)

// gitCommand is an internal reassignment of exec.Command so we
// can mock it out during testing.
var gitCommand = exec.Command

// Commit performs a git commit with a message.
func Commit(message string) (string, error) {
	cmd := gitCommand("git", "commit", "-m", message)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Add stages all files.
func Add() error {
	cmd := gitCommand("git", "add", "-A")
	return cmd.Run()
}

// ListTags lists all tags in descending order (latest at the top).
func ListTags() (string, error) {
	cmd := gitCommand("git", "tag", "--sort=-version:refname")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// LatestTag returns the name of the latest tag.
func LatestTag() (string, error) {
	cmd := gitCommand("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CreateTag creates an annotated git tag with an optional message
// if the message is an empty string, the tag name will be used.
func CreateTag(tag, message string) (string, error) {
	if message == "" {
		message = tag
	}
	cmd := gitCommand("git", "tag", "-a", tag, "-m", message)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
