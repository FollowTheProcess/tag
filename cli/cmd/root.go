// Package cmd implements the tag CLI
package cmd

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/tag/cli/app"
	"github.com/FollowTheProcess/tag/config"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version     = "dev"                                      // tag version, set at compile time by ldflags
	commit      = ""                                         // tag commit hash, set at compile time by ldflags
	headerStyle = color.New(color.FgWhite, color.Bold)       // Setting header style to use in usage message (usage.go)
	tagApp      = app.New(os.Stdout, os.Stderr, config.Path) // The tag app instance, initialised once and shared between all files in this pkg
)

// BuildRootCmd builds and returns the root tag CLI command.
func BuildRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "tag <subcommand> [flags]",
		Version:       version,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "Easy semantic versioning from the command line! 🏷",
		Long: heredoc.Doc(`
		
		Easy semantic versioning from the command line! 🏷

		`),
		Example: heredoc.Doc(`

		# See all semver tags in order
		$ tag

		# Bump a semantic version
		$ tag [patch | minor | major]
		`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return tagApp.EnsureGitRepo()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return tagApp.List()
		},
	}

	// Set our custom version and usage templates
	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetVersionTemplate(fmt.Sprintf(`{{printf "%s %s\n%s %s\n"}}`, headerStyle.Sprint("Version:"), version, headerStyle.Sprint("Commit:"), commit))

	// Add subcommands
	rootCmd.AddCommand(
		buildPatchCmd(),
		buildMinorCmd(),
		buildMajorCmd(),
		buildLatestCmd(),
	)

	return rootCmd
}
