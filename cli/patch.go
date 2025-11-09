package cli

import (
	"os"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/tag/app"
)

const (
	patchLong = `
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

// buildPatch builds and returns the patch subcommand.
func buildPatch() (*cli.Command, error) {
	var (
		push   bool
		force  bool
		dryRun bool
	)
	cmd, err := cli.New(
		"patch",
		cli.Short("Bump the patch version and issue a new tag"),
		cli.Long(patchLong),
		cli.Example("Bump the patch version", "tag patch"),
		cli.Example("Bump and push the tag to the remote", "tag patch --push"),
		cli.Example("Do not prompt for confirmation", "tag patch --push --force"),
		cli.Flag(&push, "push", 'p', false, "Push the tag to the remote"),
		cli.Flag(&force, "force", 'f', false, "Bypass confirmation prompt"),
		cli.Flag(&dryRun, "dry-run", 'd', false, "Print what would have happened"),
		cli.Run(func(cmd *cli.Command) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.Patch(push, force, dryRun)
		}),
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
