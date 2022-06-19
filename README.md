<p align="center">
<img src="https://github.com/FollowTheProcess/tag/raw/main/img/logo.png" alt="logo" width=50% height=50%>
</p>

# tag

[![License](https://img.shields.io/github/license/FollowTheProcess/tag)](https://github.com/FollowTheProcess/tag)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/tag)](https://goreportcard.com/report/github.com/FollowTheProcess/tag)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/tag?logo=github&sort=semver)](https://github.com/FollowTheProcess/tag)
[![CI](https://github.com/FollowTheProcess/tag/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/tag/actions?query=workflow%3ACI)

Easy semantic versioning from the command line! 🏷

* Free software: Apache Software License 2.0

## Project Description

Releasing new versions of software can be hard! Most projects have CI/CD pipelines set up to help with this and these pipelines are typically triggered on push of a new [semver] tag e.g. `v1.2.4`.

I made tag because I can never remember the commands to correctly issue and push a tag:

* "Was it `git tag v1.2.4`?"
* "Do I need to annotate it: `git tag -a v1.2.4`?"
* "Do I need to add a message: `git tag -a v1.2.4 -m "Some message"`?
* "Wait how do I push it again: `git push --tags` or `git push origin v1.2.4`?"

This invariably ends up with me doing it differently across every project, spending (even more) time on stackoverflow googling random git commands.

And not to mention having to replace versions in documentation, project metadata files etc.

No more 🚀 `tag` has you covered!

`tag` is cross-platform and is tested on mac, windows and linux. It'll run anywhere you can run Go!

**Fun fact:** `tag` actually releases itself!

## Installation

There are compiled executables for mac, linux and windows in the GitHub releases section, just download the correct one for your system and architecture.

There is also a [homebrew] tap:

```shell
brew tap FollowTheProcess/homebrew-tap

brew install FollowTheProcess/homebrew-tap/tag
```

## Contributing

### Developing

`tag` is a very simple project and the goal of the project is to remain very simple in line with the good old unix philosophy:

> Write programs that do one thing and do it well.
>
> Ken Thompson

Contributions are very much welcomed but please keep this goal in mind :dart:

`tag` is run as a fairly standard Go project:

* We use all the standard go tools (go test, go build etc.)
* Linting is done with the help of [golangci-lint] (see docs for install help)

We use [just] as the command runner (mainly because makefiles make me ill, but also because it's great!)

### Collaborating

No hard and fast rules here but a few guidelines:

* Raise an issue before doing a load of work on a PR, saves everyone bother
* If you add a feature, be sure to add tests to cover what you've added
* If you fix a bug, add a test that would have caught the bug you just squashed
* Be nice :smiley:

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[semver]: https://semver.org
[homebrew]: https://brew.sh
[golangci-lint]: https://golangci-lint.run
[just]: https://github.com/casey/just
[copier]: https://copier.readthedocs.io/en/latest/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier

## Regex

```go
var semVerRegex = regexp.MustCompile(fmt.Sprintf(`^v?(?P<%s>0|[1-9]\d*)\.(?P<%s>0|[1-9]\d*)\.(?P<%s>0|[1-9]\d*)(?:-(?P<%s>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<%s>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`, major, minor, patch, prerelease, buildmetadata))
```
