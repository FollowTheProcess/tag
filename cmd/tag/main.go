package main

import (
	"os"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/tag/cli"
)

func main() {
	cmd := cli.Build()
	if err := cmd.Execute(); err != nil {
		msg.Error("%s", err)
		os.Exit(1)
	}
}
