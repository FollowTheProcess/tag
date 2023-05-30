// Package cli implements tags command line interface.
package cli

import "github.com/spf13/cobra"

const (
	about   = "The all in one semver management tool 🛠️"
	example = `
# List tags in order
$ tag list

# Get latest tag
$ tag latest

# Bump a version (including content search and replace)
$ tag {patch | minor | major}
`
)

// These are all set at compile time.
var (
	version   = "dev"
	commit    = ""
	buildDate = ""
	builtBy   = ""
)

// Build builds and returns the tag CLI.
func Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tag COMMAND [FLAGS]",
		Version:       version,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         about,
		Example:       example,
	}

	// Set our custom version and usage templates
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SetVersionTemplate(versionTemplate)

	// Attach the subcommands
	cmd.AddCommand(
		buildList(),
		buildLatest(),
		buildInit(),
		buildMajor(),
		buildMinor(),
		buildPatch(),
	)

	return cmd
}
