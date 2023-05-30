package cli

import (
	"os"

	"github.com/FollowTheProcess/tag/app"
	"github.com/spf13/cobra"
)

const latestExample = `
$ tag latest
`

// buildLatest builds and returns the latest subcommand.
func buildLatest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "latest",
		Args:    cobra.NoArgs,
		Short:   "Show latest semver tag",
		Example: latestExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.Latest()
		},
	}

	return cmd
}
