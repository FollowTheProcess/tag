<p align="center">
<img src="https://github.com/FollowTheProcess/tag/raw/main/img/logo.png" alt="logo" width=50% height=50%>
</p>

# Tag

[![License](https://img.shields.io/github/license/FollowTheProcess/tag)](https://github.com/FollowTheProcess/tag)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/tag)](https://goreportcard.com/report/github.com/FollowTheProcess/tag)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/tag?logo=github&sort=semver)](https://github.com/FollowTheProcess/tag)
[![CI](https://github.com/FollowTheProcess/tag/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/tag/actions?query=workflow%3ACI)

The all in one semver management tool

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

Compiled binaries for all supported platforms can be found in the [GitHub release]. There is also a [homebrew] tap:

```shell
brew install --cask FollowTheProcess/tap/tag
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
version = '0.1.0'

[[file]]
path = 'README.md'
search = 'My project, version {{.Current}}'

[[file]]
path = 'somewhereelse.txt'
search = 'Replace me, version {{.Current}}'
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

## Config File

As mentioned above, `tag` has an optional config file (`.tag.toml`) to be placed at the root of your repo, we've seen specifying files to search and replace
contents on, but it can do a bit more than that!

A fully populated config file looks like this:

```toml
version = '0.1.0'

[git]
default-branch = 'main'
message-template = 'Bump version {{.Current}} -> {{.Next}}'
tag-template = 'v{{.Next}}'

[hooks]
pre-replace = "echo 'I run before doing anything'"
pre-commit = "echo 'I run after replacing but before committing changes'"
pre-tag = "echo 'I run after committing changes but before tagging'"
pre-push = "echo 'I run after tagging, but before pushing'"

[[file]]
path = 'pyproject.toml'
search = 'version = "{{.Current}}"'

[[file]]
path = 'README.md'
search = 'My project, version {{.Current}}'
```

### Git

The git section allows you to specify how tag interacts with git whilst bumping versions. You can specify:

* The default branch for your repo (defaults to `main`). This will be checked prior to bumping to ensure you don't issue a tag on a different branch
* The commit message template (defaults to `Bump version {{.Current}} -> {{.Next}}`). This sets the message used for your bump commit after contents have been replaced
* The tag message template (defaults to `v{{.Next}}`). Similar to the commit message but this one is associated to the tag itself.

### Hooks

Tag also lets you hook into various stages of the replacement/bumping process and inject custom logic in the form of hooks. Hooks are small shell commands that
let you update things that tag cannot see or run custom commands.

A good use case is for example, issuing a new version of a rust project with a `Cargo.toml`. In the `Cargo.toml` you must specify a version of your crate:

```toml
# Cargo.toml
version = "0.1.0"
```

When you compile your crate, it generates a `Cargo.lock` which *also* has the version. So if you use tag to bump the version in the `Cargo.toml` then the `Cargo.lock` can fall out of sync and then your crate will fail to build. Because we should never really interact with `Cargo.lock` manually, we can use hooks to re-build the crate
after replacing the version in `Cargo.toml`:

```toml
# .tag.toml
[hooks]
pre-commit = "cargo build" # Update the lockfile
```

The hooks are split into stages:

* **`pre-replace`**: This one runs first, more or less before tag does *anything* at all
* **`pre-commit`**: Runs after replacing contents, but before those changes are added and committed to the repo
* **`pre-tag`**: Runs after replacing and the changes have been committed, but before the new tag is created
* **`pre-push`**: Runs last, after everything above is finished but before the tag is pushed to the remote (if the `--push` flag is used)

[GitHub release]: https://github.com/FollowTheProcess/tag/releases
[homebrew]: https://brew.sh
[semver]: https://semver.org
