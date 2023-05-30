// Package hooks implements the hook running functionality in tag.
package hooks

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// HookStage represents a known hook stage in tag.
type HookStage int

//go:generate stringer -type=HookStage -linecomment -output=stage.go
const (
	StagePreReplace HookStage = iota // Pre Replace
	StagePreCommit                   // Pre Commit
	StagePreTag                      // Pre Tag
	StagePrePush                     // Pre Push
)

const execTimeout = 10 * time.Second

// Run runs the hook defined for a given stage, if there is no hook
// for the given stage, this becomes a no-op.
func Run(stage HookStage, cmd string, stdout, stderr io.Writer) error {
	// No-op if there's no cmd, i.e. the hook is not defined
	if cmd == "" {
		return nil
	}
	prog, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return fmt.Errorf("Command %q in hook stage %s not valid shell syntax: %w", cmd, stage, err)
	}

	runner, err := interp.New(
		interp.Params("-e"),
		interp.ExecHandler(interp.DefaultExecHandler(execTimeout)),
		interp.OpenHandler(interp.DefaultOpenHandler()),
		interp.StdIO(nil, stdout, stderr),
	)
	if err != nil {
		return fmt.Errorf("could not configure sh interpreter: %w", err)
	}

	err = runner.Run(context.Background(), prog)
	if err != nil {
		return fmt.Errorf("command %q in hook stage %s resulted in an error: %w", cmd, stage, err)
	}

	return nil
}
