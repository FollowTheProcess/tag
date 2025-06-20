package cli

import (
	"os"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/tag/app"
)

const (
	minorLong = `
You may also push the tag to any configured remote
with the "-p/--push" flag.

You will be prompted for confirmation before bumping, this
can be bypassed by passing the "-f/--force" flag.

The message accompanying the tag defaults to the tag version
itself (e.g. "v1.2.4").

If the "-d/--dry-run" flag is used, tag will simply print what would
have happened, but not do anything. This is useful for checking you have
set everything up correctly.
`
)

// buildMinor builds and returns the minor subcommand.
func buildMinor() (*cli.Command, error) {
	var (
		push   bool
		force  bool
		dryRun bool
	)
	cmd, err := cli.New(
		"minor",
		cli.Short("Bump the minor version and issue a new tag"),
		cli.Long(minorLong),
		cli.Example("Bump the minor version", "tag minor"),
		cli.Example("Bump and push the tag to the remote", "tag minor --push"),
		cli.Example("Do not prompt for confirmation", "tag minor --push --force"),
		cli.Allow(cli.NoArgs()),
		cli.Flag(&push, "push", 'p', false, "Push the tag to the remote"),
		cli.Flag(&force, "force", 'f', false, "Bypass confirmation prompt"),
		cli.Flag(&dryRun, "dry-run", 'd', false, "Print what would have happened"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.Minor(push, force, dryRun)
		}),
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
