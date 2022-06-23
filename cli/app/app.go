// Package app implements the CLI functionality, the CLI defers
// execution to the exported methods in this package
package app

import (
	"fmt"
	"io"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/tag/git"
)

// App represents the tag program.
type App struct {
	Out     io.Writer    // Where to write output to
	Printer *msg.Printer // The app's printer
}

// New creates and returns a new app.
func New(out io.Writer) *App {
	printer := msg.Default()
	printer.Out = out
	app := &App{Out: out, Printer: printer}
	return app
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
	list, err := git.ListTags()
	if err != nil {
		return err
	}
	fmt.Fprint(a.Out, list)
	return nil
}

// Latest is the tag latest subcommand.
func (a *App) Latest() error {
	latest, err := git.LatestTag()
	if err != nil {
		return err
	}
	fmt.Fprint(a.Out, latest)
	return nil
}
