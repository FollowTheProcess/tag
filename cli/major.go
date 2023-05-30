package cli

import (
	"os"

	"github.com/FollowTheProcess/tag/app"
	"github.com/spf13/cobra"
)

const (
	majorLong = `
Bump the major version

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
	majorExample = `
$ tag major

$ tag major --push

$ tag major --push --force
`
)

// buildMajor builds and returns the major subcommand.
func buildMajor() *cobra.Command {
	var (
		push   bool
		force  bool
		dryRun bool
	)
	cmd := &cobra.Command{
		Use:     "major",
		Args:    cobra.NoArgs,
		Short:   "Bump the major version and issue a new tag",
		Long:    majorLong,
		Example: majorExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tag, err := app.New(cwd, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return tag.Major(push, force, dryRun)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&push, "push", "p", false, "Push the tag to the remote")
	flags.BoolVarP(&force, "force", "f", false, "Bypass confirmation prompt")
	flags.BoolVarP(&dryRun, "dry-run", "d", false, "Print what would have happened")

	return cmd
}
