package cmd

import (
	"os"

	"github.com/FollowTheProcess/tag/cli/app"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// buildMajorCmd builds and returns the tag major CLI subcommand.
func buildMajorCmd() *cobra.Command {
	tag := app.New(os.Stdout)
	var (
		force   bool
		push    bool
		message string
	)
	majorCmd := &cobra.Command{
		Use:   "major [flags]",
		Args:  cobra.NoArgs,
		Short: "Bump the major version",
		Long: heredoc.Doc(`
		
		Bump the major version.

		You may also push the tag on bumping with
		the "-p/--push" flag.

		You will be prompted for confirmation prior to
		bumping, this can be overridden with the "-f/--force" flag.

		The message defaults to the tag itself (e.g. v1.2.4), a custom
		message can be passed using the "-m/--message" flag.
		`),
		Example: heredoc.Doc(`
		
		$ tag major

		$ tag major --push

		$ tag major --push --force

		$ tag major --message "my custom tag message"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.Major(force, push, message)
		},
	}

	flags := majorCmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Bypass confirmation prompt")
	flags.BoolVarP(&push, "push", "p", false, "Push tag to configured remote")
	flags.StringVarP(&message, "message", "m", "", "Custom tag message.")

	return majorCmd
}
