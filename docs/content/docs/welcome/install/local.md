---
title: "Local Installation"
weight: 2
---

## Binaries

```bash
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin {{< golangci/latest-version >}}

# or install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s {{< golangci/latest-version >}}

# In Alpine Linux (as it does not come with curl by default)
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s {{< golangci/latest-version >}}

golangci-lint --version
```

On Windows, you can run the above commands with Git Bash, which comes with [Git for Windows](https://git-scm.com/download/win).

## Linux

Golangci-lint is available inside the majority of the package managers.

{{% details closed="true" title="Packaging status" %}}

[![Packaging status](https://repology.org/badge/vertical-allrepos/golangci-lint.svg)](https://repology.org/project/golangci-lint/versions)

{{% /details %}}

## macOS

### Homebrew

Note: Homebrew can use an unexpected version of Go to build the binary,
so we recommend either using our binaries or ensuring the version of Go used to build.

You can install a binary release on macOS using [brew](https://brew.sh/):

```bash
brew install golangci-lint
brew upgrade golangci-lint
```

Note: Previously, we used a [Homebrew tap](https://github.com/golangci/homebrew-tap).
We recommend using the [official formula](https://formulae.brew.sh/formula/golangci-lint) instead of the tap,
but sometimes the most recent release isn't immediately available via Homebrew core due to manual updates that need to occur from Homebrew core maintainers.
In this case, the tap formula, which is updated automatically,
can be used to install the latest version of golangci-lint:

```bash
brew tap golangci/tap
brew install golangci/tap/golangci-lint
```

### MacPorts

It can also be installed through [MacPorts](https://www.macports.org/)
The MacPorts installation mode is community-driven and not officially maintained by the golangci team.

```bash
sudo port install golangci-lint
```

## Windows

### Chocolatey

You can install a binary on Windows using [chocolatey](https://community.chocolatey.org/packages/golangci-lint).

```bash
choco install golangci-lint
```

### Scoop

You can install a binary on Windows using [scoop](https://scoop.sh).

```bash
scoop install main/golangci-lint
```

The scoop package is not officially maintained by golangci team.

## Docker

The Docker image is available on [Docker Hub](https://hub.docker.com/r/golangci/golangci-lint).

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:{{< golangci/latest-version >}} golangci-lint run
```

Colored output:
```bash
docker run -t --rm -v $(pwd):/app -w /app golangci/golangci-lint:{{< golangci/latest-version >}} golangci-lint run
```

Preserving caches between consecutive runs:
```bash
docker run --rm -t -v $(pwd):/app -w /app \
--user $(id -u):$(id -g) \
-v $(go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
-v $(go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
-v ~/.cache/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
golangci/golangci-lint:{{< golangci/latest-version >}} golangci-lint run
```

## mise

Note: `mise` is using the [aqua](https://aquaproj.github.io/) backend for this tool, so binaries installed came from GitHub assets (recommended).

You can install golangci-lint by using [`mise`](https://github.com/jdx/mise).

```bash
mise use -g golangci-lint@{{< golangci/latest-version >}}
```

The `mise` integration is not officially maintained by golangci team.

## Install from Sources

> [!WARNING]
> Using `go install`/`go get`, "tools pattern", and `tool` command/directives installations aren't guaranteed to work.  
> We recommend using binary installation.

These installations aren't recommended because of the following points:

1. These installations compile golangci-lint locally. The Go version used to build will depend on your local Go version.
2. Some users use the `-u` flag for `go get`, which upgrades our dependencies. The resulting binary was not tested and is not guaranteed to work.
3. When using the "tools pattern" or `tool` command/directives, the dependencies of a tool can modify the dependencies of another tool or your project. The resulting binary was not tested and is not guaranteed to work.
4. We've encountered issues with Go module hashes due to the unexpected recreation of dependency tags.
5. `go.mod` replacement directives don't apply transitively. It means a user will be using a patched version of golangci-lint if we use such replacements.
6. It allows installation from the main branch, which can't be considered stable.
7. It's slower than binary installation.

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{< golangci/latest-version >}}
```

{{% details title="`go tool` usage recommendations" closed="true" %}}

> [!WARNING]
> We don't recommend using `go tool`.

But if you want to use `go tool` to install and run golangci-lint (**once again we don't recommend that**),
the best approach is to use a dedicated module or module file to isolate golangci-lint from other tools or dependencies.

This approach avoids modifying your project dependencies and the golangci-lint dependencies.

> [!CAUTION]
> You should never update golangci-lint dependencies manually.

**Method 1: dedicated module file**

```sh
# Create a dedicated module file
go mod init -modfile=golangci-lint.mod <your_module_path>/golangci-lint
# Example: go mod init -modfile=golangci-lint.mod github.com/org/repo/golangci-lint
```

```sh
# Add golangci-lint as a tool
go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{< golangci/latest-version >}}
```

```sh
# Run golangci-lint as a tool
go tool -modfile=golangci-lint.mod golangci-lint run
```

```sh
# Update golangci-lint
go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

**Method 2: dedicated module**

```sh
# Create a dedicated directory
mkdir golangci-lint
```

```sh
# Create a dedicated module file
go mod init -modfile=tools/go.mod <your_module_path>/golangci-lint
# Example: go mod init -modfile=golangci-lint/go.mod github.com/org/repo/golangci-lint
```

```sh
# Setup a Go workspace
go work init . golangci-lint
```

```sh
# Add golangci-lint as a tool
go get -tool -modfile=golangci-lint/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint
```

```sh
# Run golangci-lint as a tool
go tool golangci-lint run
```

```sh
# Update golangci-lint
go get -tool -modfile=golangci-lint/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

{{% /details %}}
