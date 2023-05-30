package cli

import (
	"os"

	"github.com/FollowTheProcess/tag/app"
	"github.com/spf13/cobra"
)

const (
	defaultLimit = 10
	listExample  = `
$ tag list

$ tag list --limit 15`
)

// buildList builds and returns the list subcommand.
func buildList() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:     "list",
		Args:    cobra.NoArgs,
		Short:   "Show semver tags in order",
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.List(limit)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&limit, "limit", "l", defaultLimit, "Max number of tags to show")

	return cmd
}
