package main

import (
	"context"
	"fmt"
	"os"

	"go.followtheprocess.codes/msg"
	"go.followtheprocess.codes/tag/cli"
)

func main() {
	if err := run(); err != nil {
		msg.Error("%s", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cmd, err := cli.Build()
	if err != nil {
		return fmt.Errorf("could not build tag cli: %w", err)
	}

	return cmd.Execute(ctx)
}
