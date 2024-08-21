// Package cli implements tags command line interface.
package cli

import "github.com/FollowTheProcess/cli"

// These are all set at compile time.
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// Build builds and returns the tag CLI.
func Build() (*cli.Command, error) {
	initCmd, err := buildInit()
	if err != nil {
		return nil, err
	}
	latestCmd, err := buildLatest()
	if err != nil {
		return nil, err
	}
	listCmd, err := buildList()
	if err != nil {
		return nil, err
	}
	majorCmd, err := buildMajor()
	if err != nil {
		return nil, err
	}
	minorCmd, err := buildMinor()
	if err != nil {
		return nil, err
	}
	patchCmd, err := buildPatch()
	if err != nil {
		return nil, err
	}
	cmd, err := cli.New(
		"tag",
		cli.Short("The all in one semver management tool 🛠️"),
		cli.Example("List tags in order", "tag list"),
		cli.Example("Get latest tag", "tag latest"),
		cli.Example("Bump a version (including content search and replace)", "tag {patch | minor | major}"),
		cli.Version(version),
		cli.Commit(commit),
		cli.BuildDate(buildDate),
		cli.Allow(cli.NoArgs()),
		cli.SubCommands(
			initCmd,
			latestCmd,
			listCmd,
			majorCmd,
			minorCmd,
			patchCmd,
		),
	)
	// TODO(@FollowTheProcess): Subcommands
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
