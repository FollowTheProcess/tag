package main

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/tag/cli"
)

func main() {
	if err := run(); err != nil {
		msg.Error("%s", err)
		os.Exit(1)
	}
}

func run() error {
	cmd, err := cli.Build()
	if err != nil {
		return fmt.Errorf("could not build tag cli: %w", err)
	}
	return cmd.Execute()
}
