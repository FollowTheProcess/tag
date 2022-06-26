// Package git implements tags interface with git in order to interact
// with tags and make commits.
package git

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
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

// Push performs a git push to the configured remote.
func Push() (string, error) {
	cmd := gitCommand("git", "push")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

// PushTag pushes a certain tag to the configured remote.
func PushTag(tag string) (string, error) {
	cmd := gitCommand("git", "push", "origin", tag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

// ListTags lists all tags in descending order (latest at the top).
func ListTags() (string, error) {
	// git will return nothing if there are no tags
	cmd := gitCommand("git", "tag", "--sort=-version:refname")
	out, err := cmd.CombinedOutput()
	if bytes.Equal(out, []byte("")) {
		return "", errors.New("No tags found")
	}
	return string(out), err
}

// LatestTag returns the name of the latest tag.
func LatestTag() (string, error) {
	cmd := gitCommand("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.CombinedOutput()
	if bytes.Contains(out, []byte("fatal: No names found")) {
		return "", errors.New("No tags found")
	}
	return strings.TrimSpace(string(out)), err
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

// IsRepo detects whether or not we are currently in a git repo.
func IsRepo() bool {
	cmd := gitCommand("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	if !bytes.Equal(bytes.TrimSpace(bytes.ToLower(out)), []byte("true")) {
		return false
	}
	return true
}

// Branch gets the name of the current git branch.
func Branch() (string, error) {
	cmd := gitCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsDirty checks whether or not the working tree is dirty.
func IsDirty() (bool, error) {
	cmd := gitCommand("git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	status := strings.TrimSpace(string(out))
	if len(status) > 1 {
		return true, nil
	}
	return false, nil
}
