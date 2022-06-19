package cmd

import (
	"os"

	"github.com/FollowTheProcess/tag/cli/app"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// buildPatchCmd builds and returns the tag patch CLI subcommand.
func buildPatchCmd() *cobra.Command {
	tag := &app.App{Out: os.Stdout}
	var (
		force   bool
		push    bool
		message string
	)
	patchCmd := &cobra.Command{
		Use:   "patch [flags]",
		Args:  cobra.NoArgs,
		Short: "Bump the patch version",
		Long: heredoc.Doc(`
		
		Bump the patch version.

		You may also push the tag on bumping with
		the "-p/--push" flag.

		You will be prompted for confirmation prior to
		bumping, this can be overridden with the "-f/--force" flag.

		The message defaults to the tag itself (e.g. v1.2.4), a custom
		message can be passed using the "-m/--message" flag.
		`),
		Example: heredoc.Doc(`
		
		$ tag patch

		$ tag patch --push

		$ tag patch --push --force

		$ tag patch --message "my custom tag message"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.Patch(force, push, message)
		},
	}

	flags := patchCmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Bypass confirmation prompt")
	flags.BoolVarP(&push, "push", "p", false, "Push tag to configured remote")
	flags.StringVarP(&message, "message", "m", "", "Custom tag message.")

	return patchCmd
}
