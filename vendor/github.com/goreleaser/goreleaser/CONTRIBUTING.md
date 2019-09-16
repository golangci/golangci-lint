# Contributing

By participating to this project, you agree to abide our [code of
conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`goreleaser` is written in [Go](https://golang.org/).

Prerequisites:

- `make`
- [Go 1.11+](https://golang.org/doc/install)
- `rpmbuild` (`apt get install rpm`/`brew install rpm`)
- [snapcraft](https://snapcraft.io/)
- [Docker](https://www.docker.com/)
- `gpg` (probably already installed on your system)

Clone `goreleaser` anywhere:

```sh
$ git clone git@github.com:goreleaser/goreleaser.git
```

Install the build and lint dependencies:

```sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

```sh
$ make test
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

```sh
$ make build
```

When you are satisfied with the changes, we suggest you run:

```sh
$ make ci
```

Which runs all the linters and tests.

## Create a commit

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Submit a pull request

Push your branch to your `goreleaser` fork and open a pull request against the
master branch.

## Financial contributions

We also welcome financial contributions in full transparency on our [open collective](https://opencollective.com/goreleaser).
Anyone can file an expense. If the expense makes sense for the development of the community, it will be "merged" in the ledger of our open collective by the core contributors and the person who filed the expense will be reimbursed.

## Credits

### Contributors

Thank you to all the people who have already contributed to GoReleaser!
<a href="graphs/contributors"><img src="https://opencollective.com/goreleaser/contributors.svg?width=890" /></a>

### Backers

Thank you to all our backers! [[Become a backer](https://opencollective.com/goreleaser#backer)]

<a href="https://opencollective.com/goreleaser#backers" target="_blank"><img src="https://opencollective.com/goreleaser/backers.svg?width=890"></a>

### Sponsors

Thank you to all our sponsors! (please ask your company to also support this open source project by [becoming a sponsor](https://opencollective.com/goreleaser#sponsor))

<a href="https://opencollective.com/goreleaser/sponsor/0/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/1/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/2/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/3/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/4/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/5/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/6/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/7/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/8/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/goreleaser/sponsor/9/website" target="_blank"><img src="https://opencollective.com/goreleaser/sponsor/9/avatar.svg"></a>
