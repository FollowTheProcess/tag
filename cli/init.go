package cli

import (
	"os"

	"github.com/FollowTheProcess/tag/app"
	"github.com/spf13/cobra"
)

const initExample = `
$ tag init

$ tag init --force
`

// buildInit builds and returns the init subcommand.
func buildInit() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:     "init",
		Args:    cobra.NoArgs,
		Short:   "Create a new tag config file",
		Example: initExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}

			return tag.Init(cwd, force)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Overwrite an existing config file")

	return cmd
}
