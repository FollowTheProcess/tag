package cli

import (
	"os"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/tag/app"
)

// buildLatest builds and returns the latest subcommand.
func buildLatest() (*cli.Command, error) {
	cmd, err := cli.New(
		"latest",
		cli.Short("Show latest semver tag"),
		cli.Example("Show the latest", "tag latest"),
		cli.Allow(cli.NoArgs()),
		cli.Run(func(cmd *cli.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.Latest()
		}),
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
