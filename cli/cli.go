// Package cli implements tags command line interface.
package cli

import "go.followtheprocess.codes/cli"

// These are all set at compile time.
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// Build builds and returns the tag CLI.
func Build() (*cli.Command, error) {
	cmd, err := cli.New(
		"tag",
		cli.Short("The all in one semver management tool 🛠️"),
		cli.Example("List tags in order", "tag list"),
		cli.Example("Get latest tag", "tag latest"),
		cli.Example("Bump a version (including content search and replace)", "tag {patch | minor | major}"),
		cli.Version(version),
		cli.Commit(commit),
		cli.BuildDate(buildDate),
		cli.SubCommands(
			buildInit,
			buildLatest,
			buildList,
			buildMajor,
			buildMinor,
			buildPatch,
		),
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
