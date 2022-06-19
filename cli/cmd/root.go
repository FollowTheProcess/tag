// Package cmd implements the tag CLI
package cmd

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/tag/cli/app"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version     = "dev"                                // tag version, set at compile time by ldflags
	commit      = ""                                   // tag commit hash, set at compile time by ldflags
	headerStyle = color.New(color.FgWhite, color.Bold) // Setting header style to use in usage message (usage.go)
)

// BuildRootCmd builds and returns the root tag CLI command.
func BuildRootCmd() *cobra.Command {
	// Note: options must be a pointer so flags are propegated to the App struct
	options := &app.Options{}
	tag := &app.App{
		Out:     os.Stdout,
		Options: options,
	}

	rootCmd := &cobra.Command{
		Use:           "tag [arguments]...",
		Version:       version,
		Args:          cobra.ArbitraryArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "Easy semantic versioning from the command line",
		Long: heredoc.Doc(`
		
		Easy semantic versioning from the command line

		`),
		Example: heredoc.Doc(`

		# Some examples of your CLI
		$ tag --help
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.Run(args)
		},
	}

	// Attach the flags
	flags := rootCmd.Flags()
	flags.BoolVar(&options.Switch, "switch", false, "Some boolean switch")

	// Set our custom version and usage templates
	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetVersionTemplate(fmt.Sprintf(`{{printf "%s %s\n%s %s\n"}}`, headerStyle.Sprint("Version:"), version, headerStyle.Sprint("Commit:"), commit))

	return rootCmd
}
