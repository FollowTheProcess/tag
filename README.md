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
brew install FollowTheProcess/homebrew-tap/tag
```

## Usage

Tag has 2 modes of operating, one in which it doesn't find a config file in the current directory (`.tag.toml`), and one where it does. Let's start with the first mode.

### No Replace Mode

If there is no config file present in `cwd`, tag will operate in "no replace" mode. This is it's most basic mode and when tag
is in this mode all you can do with it is list, create, and push new [semver] tags.

For example let's say you're working on a project currently at `v0.23.8` and you've decided you want to signal to the world that your project is stable, it's time for a major version bump! 🚀

Your project also has a CI/CD pipeline where on the push of a new tag it gets compiled and packaged up and a new release gets created.

So you need to create a new tag (`v1.0.0`) and push it. No problem!

```shell
tag major --push
```

This will create a new `v1.0.0` annotated git tag, and push it to the configured remote. Job done ✅

### Replace Mode

Now this is already nice but wouldn't it be *even nicer* if you didn't have to manually bump version numbers in project metadata files, or maybe the README:

```markdown
# My Project Readme

This my project, version = 0.1.0
```

`tag` can do that too! All you have to do is tell it what to do with which files to work on, enter the `.tag.toml` config file which should be placed in the root of your repo:

```toml
[tag]
files = [
    { path = "README.md", search = "version = {{.Current}}", replace = "version = {{.Next}}" },
]
```

Tag uses two special variables `{{.Current}}` and `{{.Next}}` to substitute for the correct versions while bumping as well as the path (relative to `.tag.toml`) of the files you want to change.

So now all you have to do is e.g.

```shell
tag minor --push
```

And then tag will:

* Perform search and replace on all occurrences of your search string
* Stage all the changes in git once the replacing is done
* Commit the changes with a message like `Bump version 0.1.0 -> 0.2.0`
* Push the changes
* Push the new tag

And then your CI/CD pipeline will take care of the rest! 🎉

After bumping, your README will now look like this:

```markdown
# My Project Readme

This my project, version = 0.2.0
```

## Contributing

### Developing

`tag` is a very simple project and the goal of the project is to remain very simple in line with the good old unix philosophy:

> Write programs that do one thing and do it well.
>
> Ken Thompson

Contributions are very much welcomed but please keep this goal in mind 🎯

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
