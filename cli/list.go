package cli

import (
	"context"
	"os"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/tag/app"
)

const (
	defaultLimit = 10
)

// buildList builds and returns the list subcommand.
func buildList() (*cli.Command, error) {
	var limit int
	cmd, err := cli.New(
		"list",
		cli.Short("Show semver tags in order"),
		cli.Example("Show all tags", "tag list"),
		cli.Example("Limit to a max number", "tag list --limit 15"),
		cli.Flag(&limit, "limit", 'l', "Max number of tags to show", cli.FlagDefault(defaultLimit)),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.List(limit)
		}),
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
