// Package app implements the CLI functionality, the CLI defers
// execution to the exported methods in this package
package app

import (
	"fmt"
	"io"
)

// App represents the tag program.
type App struct {
	Out     io.Writer
	Options *Options
}

// Options holds all the flag options for tag, these will be at their zero values
// if the flags were not set and the value of the flag otherwise.
type Options struct {
	Switch bool // Some boolean switch
}

// Run is the entry point to the tag program.
func (a *App) Run(args []string) error {
	fmt.Fprintf(a.Out, "Args: %v\n", args)
	fmt.Fprintf(a.Out, "Switch flag was %v\n", a.Options.Switch)

	return nil
}
