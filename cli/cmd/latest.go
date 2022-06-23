package cmd

import (
	"os"

	"github.com/FollowTheProcess/tag/cli/app"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// buildLatestCmd builds and returns the tag latest CLI subcommand.
func buildLatestCmd() *cobra.Command {
	tag := app.New(os.Stdout)
	latestCmd := &cobra.Command{
		Use:   "latest",
		Args:  cobra.NoArgs,
		Short: "Show the latest semver tag",
		Example: heredoc.Doc(`
		
		$ tag latest
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.Latest()
		},
	}

	return latestCmd
}
