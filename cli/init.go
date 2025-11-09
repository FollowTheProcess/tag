package cli

import (
	"os"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/tag/app"
)

// buildInit builds and returns the init subcommand.
func buildInit() (*cli.Command, error) {
	var force bool
	cmd, err := cli.New(
		"init",
		cli.Short("Create a new tag config file"),
		cli.Example("Create a config file", "tag init"),
		cli.Example("Overwrite an existing one", "tag init --force"),
		cli.Flag(&force, "force", 'f', false, "Overwrite an existing config file"),
		cli.Run(func(cmd *cli.Command) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}

			return tag.Init(cwd, force)
		}),
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
