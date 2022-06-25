package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// buildMinorCmd builds and returns the tag minor CLI subcommand.
func buildMinorCmd() *cobra.Command {
	var (
		force   bool
		push    bool
		message string
	)
	minorCmd := &cobra.Command{
		Use:   "minor [flags]",
		Args:  cobra.NoArgs,
		Short: "Bump the minor version",
		Long: heredoc.Doc(`
		
		Bump the minor version.

		You may also push the tag on bumping with
		the "-p/--push" flag.

		You will be prompted for confirmation prior to
		bumping, this can be overridden with the "-f/--force" flag.

		The message defaults to the tag itself (e.g. v1.2.4), a custom
		message can be passed using the "-m/--message" flag.
		`),
		Example: heredoc.Doc(`
		
		$ tag minor

		$ tag minor --push

		$ tag minor --push --force

		$ tag minor --message "my custom tag message"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tagApp.Minor(force, push, message)
		},
	}

	flags := minorCmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Bypass confirmation prompt")
	flags.BoolVarP(&push, "push", "p", false, "Push tag to configured remote")
	flags.StringVarP(&message, "message", "m", "", "Custom tag message.")

	return minorCmd
}
