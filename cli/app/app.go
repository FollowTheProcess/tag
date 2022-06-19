// Package app implements the CLI functionality, the CLI defers
// execution to the exported methods in this package
package app

import (
	"fmt"
	"io"
)

// App represents the tag program.
type App struct {
	Out io.Writer // Where to write output to
}

// Patch is the tag patch subcommand.
func (a *App) Patch(force, push bool, message string) error {
	fmt.Printf("force: %v\n", force)
	fmt.Printf("push: %v\n", push)
	fmt.Printf("message: %s\n", message)
	return nil
}

// Minor is the tag minor subcommand.
func (a *App) Minor(force, push bool, message string) error {
	fmt.Printf("force: %v\n", force)
	fmt.Printf("push: %v\n", push)
	fmt.Printf("message: %s\n", message)
	return nil
}

// Major is the tag major subcommand.
func (a *App) Major(force, push bool, message string) error {
	fmt.Printf("force: %v\n", force)
	fmt.Printf("push: %v\n", push)
	fmt.Printf("message: %s\n", message)
	return nil
}

// List is what happens when tag is invoked with no subcommands.
func (a *App) List() error {
	fmt.Println("No subcommands, list all tags in order")
	return nil
}

// Latest is the tag latest subcommand.
func (a *App) Latest() error {
	fmt.Println("Get latest tag")
	return nil
}
