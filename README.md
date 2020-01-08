# GolangCI-Lint

[![Build Status](https://travis-ci.com/golangci/golangci-lint.svg?branch=master)](https://travis-ci.com/golangci/golangci-lint)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)
[![License](https://img.shields.io/github/license/golangci/golangci-lint)](/LICENSE)
[![Release](https://img.shields.io/github/release/golangci/golangci-lint.svg)](https://github.com/golangci/golangci-lint/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/golangci/golangci-lint)](https://hub.docker.com/r/golangci/golangci-lint)

GolangCI-Lint is a linters aggregator. It's fast: on average [5 times faster](#performance) than gometalinter.
It's [easy to integrate and use](#command-line-options), has [nice output](#quick-start) and has a minimum number of false positives. It supports go modules.

GolangCI-Lint has [integrations](#editor-integration) with VS Code, GNU Emacs, Sublime Text.

Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).

Sponsored by [GolangCI.com](https://golangci.com): SaaS service for running linters on GitHub pull requests. Free for Open Source.

<a href="https://golangci.com/"><img src="docs/go.png" width="250px"></a>

- [GolangCI-Lint](#golangci-lint)
  - [Demo](#demo)
  - [Install](#install)
    - [Binary](#binary)
    - [macOS](#macos)
    - [Docker](#docker)
    - [Go](#go)
  - [Trusted By](#trusted-by)
  - [Quick Start](#quick-start)
  - [Editor Integration](#editor-integration)
  - [Shell Completion](#shell-completion)
    - [macOS](#macos-1)
    - [Linux](#linux)
  - [Comparison](#comparison)
    - [`golangci-lint` vs `gometalinter`](#golangci-lint-vs-gometalinter)
    - [`golangci-lint` vs Running Linters Manually](#golangci-lint-vs-running-linters-manually)
  - [Performance](#performance)
    - [Comparison with gometalinter](#comparison-with-gometalinter)
    - [Why golangci-lint is faster](#why-golangci-lint-is-faster)
    - [Memory Usage of Golangci-lint](#memory-usage-of-golangci-lint)
  - [Internals](#internals)
  - [Supported Linters](#supported-linters)
    - [Enabled By Default Linters](#enabled-by-default-linters)
    - [Disabled By Default Linters (`-E/--enable`)](#disabled-by-default-linters--e--enable)
  - [Configuration](#configuration)
    - [Command-Line Options](#command-line-options)
    - [Config File](#config-file)
  - [False Positives](#false-positives)
    - [Nolint](#nolint)
  - [FAQ](#faq)
  - [Thanks](#thanks)
  - [Changelog](#changelog)
  - [Debug](#debug)
  - [Future Plans](#future-plans)
  - [Contact Information](#contact-information)
  - [License Scan](#license-scan)

## Demo

<p align="center">
  <img src="./docs/demo.svg" width="100%">
</p>

Short 1.5 min video demo of analyzing [beego](https://github.com/astaxie/beego).
[![asciicast](https://asciinema.org/a/183662.png)](https://asciinema.org/a/183662)

## Install

### Binary

Most installations are done for CI (e.g. Travis CI, CircleCI). It's important to have reproducible CI:
don't start to fail all builds at the same time. With golangci-lint this can happen if you
use deprecated option `--enable-all` and a new linter is added or even without `--enable-all`: when one upstream linter is upgraded.

It's highly recommended to install a specific version of golangci-lint available on the [releases page](https://github.com/golangci/golangci-lint/releases).

Here is the recommended way to install golangci-lint v1.22.2:

```bash
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.22.2

# or install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.22.2

# In alpine linux (as it does not come with curl by default)
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.22.2

golangci-lint --version
```

It is advised that you periodically update version of golangci-lint as the project is under active development
and is constantly being improved. For any problems with golangci-lint, check out recent [GitHub issues](https://github.com/golangci/golangci-lint/issues) and update if needed.

### macOS

You can also install a binary release on macOS using [brew](https://brew.sh/):

```bash
brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint
```

### Docker

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.22.2 golangci-lint run -v
```

### Go

Please, do not install `golangci-lint` by `go get`:

1. [`go.mod`](https://github.com/golangci/golangci-lint/blob/master/go.mod) replacement directive doesn't apply. It means you will be using patched version of `golangci-lint`.
2. it's much slower than binary installation
3. its stability depends on your Go version (e.g. on [this compiler Go <= 1.12 bug](https://github.com/golang/go/issues/29612)).
4. it's not guaranteed to work: e.g. we've encountered a lot of issues with Go modules hashes.
5. it allows installation from `master` branch which can't be considered stable.

## Trusted By

The following companies/products use golangci-lint:

* [Google](https://github.com/GoogleContainerTools/skaffold)
* [Facebook](https://github.com/facebookincubator/fbender)
* [Red Hat OpenShift](https://github.com/openshift/telemeter)
* [Yahoo](https://github.com/yahoo/yfuzz)
* [IBM](https://github.com/ibm-developer/ibm-cloud-env-golang)
* [Xiaomi](https://github.com/XiaoMi/soar)
* [Baidu](https://github.com/baidu/bfe)
* [Samsung](https://github.com/samsung-cnct/cluster-api-provider-ssh)
* [Arduino](https://github.com/arduino/arduino-cli)
* [Eclipse Foundation](https://github.com/eclipse/che-go-jsonrpc)
* [WooCart](https://github.com/woocart/gsutil)
* [Percona](https://github.com/percona/pmm-managed)
* [Serverless](https://github.com/serverless/event-gateway)
* [ScyllaDB](https://github.com/scylladb/gocqlx)
* [NixOS](https://github.com/NixOS/nixpkgs-channels)
* [The New York Times](https://github.com/NYTimes/encoding-wrapper)
* [Istio](https://github.com/istio/istio)
* [SoundCloud](https://github.com/soundcloud/periskop)
* [Mattermost](https://github.com/mattermost/mattermost-server)

The following great projects use golangci-lint:

* [alecthomas/participle](https://github.com/alecthomas/participle)
* [asobti/kube-monkey](https://github.com/asobti/kube-monkey)
* [banzaicloud/pipeline](https://github.com/banzaicloud/pipeline)
* [caicloud/cyclone](https://github.com/caicloud/cyclone)
* [getantibody/antibody](https://github.com/getantibody/antibody)
* [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser)
* [go-swagger/go-swagger](https://github.com/go-swagger/go-swagger)
* [kubeedge/kubeedge](https://github.com/kubeedge/kubeedge)
* [kubernetes-sigs/kustomize](https://github.com/kubernetes-sigs/kustomize)
* [dunglas/mercure](https://github.com/dunglas/mercure)
* [posener/complete](https://github.com/posener/complete)
* [segmentio/terraform-docs](https://github.com/segmentio/terraform-docs)
* [tsuru/tsuru](https://github.com/tsuru/tsuru)
* [twpayne/chezmoi](https://github.com/twpayne/chezmoi)
* [virtual-kubelet/virtual-kubelet](https://github.com/virtual-kubelet/virtual-kubelet)
* [xenolf/lego](https://github.com/xenolf/lego)
* [y0ssar1an/q](https://github.com/y0ssar1an/q)

## Quick Start

To run golangci-lint execute:

```bash
golangci-lint run
```

It's an equivalent of executing:

```bash
golangci-lint run ./...
```

You can choose which directories and files to analyze:

```bash
golangci-lint run dir1 dir2/... dir3/file1.go
```

Directories are NOT analyzed recursively. To analyze them recursively append `/...` to their path.

GolangCI-Lint can be used with zero configuration. By default the following linters are enabled:

```bash
$ golangci-lint help linters
Enabled by default linters:
deadcode: Finds unused code [fast: true, auto-fix: false]
errcheck: Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases [fast: true, auto-fix: false]
gosimple (megacheck): Linter for Go source code that specializes in simplifying a code [fast: true, auto-fix: false]
govet (vet, vetshadow): Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string [fast: true, auto-fix: false]
ineffassign: Detects when assignments to existing variables are not used [fast: true, auto-fix: false]
staticcheck (megacheck): Staticcheck is a go vet on steroids, applying a ton of static analysis checks [fast: true, auto-fix: false]
structcheck: Finds unused struct fields [fast: true, auto-fix: false]
typecheck: Like the front-end of a Go compiler, parses and type-checks Go code [fast: true, auto-fix: false]
unused (megacheck): Checks Go code for unused constants, variables, functions and types [fast: false, auto-fix: false]
varcheck: Finds unused global variables and constants [fast: true, auto-fix: false]
```

and the following linters are disabled by default:

```bash
$ golangci-lint help linters
...
Disabled by default linters:
bodyclose: checks whether HTTP response body is closed successfully [fast: true, auto-fix: false]
depguard: Go linter that checks if package imports are in a list of acceptable packages [fast: true, auto-fix: false]
dogsled: Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f()) [fast: true, auto-fix: false]
dupl: Tool for code clone detection [fast: true, auto-fix: false]
funlen: Tool for detection of long functions [fast: true, auto-fix: false]
gochecknoglobals: Checks that no globals are present in Go code [fast: true, auto-fix: false]
gochecknoinits: Checks that no init functions are present in Go code [fast: true, auto-fix: false]
gocognit: Computes and checks the cognitive complexity of functions [fast: true, auto-fix: false]
goconst: Finds repeated strings that could be replaced by a constant [fast: true, auto-fix: false]
gocritic: The most opinionated Go source code linter [fast: true, auto-fix: false]
gocyclo: Computes and checks the cyclomatic complexity of functions [fast: true, auto-fix: false]
godox: Tool for detection of FIXME, TODO and other comment keywords [fast: true, auto-fix: false]
gofmt: Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification [fast: true, auto-fix: true]
goimports: Goimports does everything that gofmt does. Additionally it checks unused imports [fast: true, auto-fix: true]
golint: Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes [fast: true, auto-fix: false]
gomnd: An analyzer to detect magic numbers. [fast: true, auto-fix: false]
gosec (gas): Inspects source code for security problems [fast: true, auto-fix: false]
interfacer: Linter that suggests narrower interface types [fast: true, auto-fix: false]
lll: Reports long lines [fast: true, auto-fix: false]
maligned: Tool to detect Go structs that would take less memory if their fields were sorted [fast: true, auto-fix: false]
misspell: Finds commonly misspelled English words in comments [fast: true, auto-fix: true]
nakedret: Finds naked returns in functions greater than a specified function length [fast: true, auto-fix: false]
prealloc: Finds slice declarations that could potentially be preallocated [fast: true, auto-fix: false]
scopelint: Scopelint checks for unpinned variables in go programs [fast: true, auto-fix: false]
stylecheck: Stylecheck is a replacement for golint [fast: true, auto-fix: false]
unconvert: Remove unnecessary type conversions [fast: true, auto-fix: false]
unparam: Reports unused function parameters [fast: true, auto-fix: false]
whitespace: Tool for detection of leading and trailing whitespace [fast: true, auto-fix: true]
wsl: Whitespace Linter - Forces you to use empty lines! [fast: true, auto-fix: false]
```

Pass `-E/--enable` to enable linter and `-D/--disable` to disable:

```bash
golangci-lint run --disable-all -E errcheck
```

## Editor Integration

1. [Go for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go).
   Recommended settings for VS Code are:

   ```json
   "go.lintTool":"golangci-lint",
   "go.lintFlags": [
     "--fast"
   ]
   ```

   Using it in an editor without `--fast` can freeze your editor.
   Golangci-lint automatically discovers `.golangci.yml` config for edited file: you don't need to configure it in VS Code settings.
2. Sublime Text - [plugin](https://github.com/alecthomas/SublimeLinter-contrib-golang-cilint) for SublimeLinter.
3. GoLand
   * Add [File Watcher](https://www.jetbrains.com/help/go/settings-tools-file-watchers.html) using existing `golangci-lint` template.
   * If your version of GoLand does not have the `golangci-lint` [File Watcher](https://www.jetbrains.com/help/go/settings-tools-file-watchers.html) template you can configure your own and use arguments `run --disable=typecheck $FileDir$`.
4. GNU Emacs
   * [Spacemacs](https://github.com/syl20bnr/spacemacs/blob/develop/layers/+lang/go/README.org#pre-requisites)
   * [flycheck checker](https://github.com/weijiangan/flycheck-golangci-lint).
5. Vim
   * [vim-go](https://github.com/fatih/vim-go)
   * syntastic [merged pull request](https://github.com/vim-syntastic/syntastic/pull/2190) with golangci-lint support
   * ale [merged pull request](https://github.com/w0rp/ale/pull/1890) with golangci-lint support
6. Atom - [go-plus](https://atom.io/packages/go-plus) supports golangci-lint.

## Shell Completion

`golangci-lint` can generate bash completion file.

### macOS

There are two versions of `bash-completion`, v1 and v2. V1 is for Bash 3.2 (which is the default on macOS), and v2 is for Bash 4.1+. The `golangci-lint` completion script doesn’t work correctly with bash-completion v1 and Bash 3.2. It requires bash-completion v2 and Bash 4.1+. Thus, to be able to correctly use `golangci-lint` completion on macOS, you have to install and use Bash 4.1+ ([instructions](https://itnext.io/upgrading-bash-on-macos-7138bd1066ba)). The following instructions assume that you use Bash 4.1+ (that is, any Bash version of 4.1 or newer).

Install `bash-completion v2`:

```bash
brew install bash-completion@2
echo 'export BASH_COMPLETION_COMPAT_DIR="/usr/local/etc/bash_completion.d"' >>~/.bashrc
echo '[[ -r "/usr/local/etc/profile.d/bash_completion.sh" ]] && . "/usr/local/etc/profile.d/bash_completion.sh"' >>~/.bashrc
exec bash # reload and replace (if it was updated) shell
type _init_completion && echo "completion is OK" # verify that bash-completion v2 is correctly installed
```

Add `golangci-lint` bash completion:

```bash
echo 'source <(golangci-lint completion bash)' >>~/.bashrc
source ~/.bashrc
```

### Linux

See [kubectl instructions](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enabling-shell-autocompletion) and don't forget to replace `kubectl` with `golangci-lint`.

## Comparison

### `golangci-lint` vs `gometalinter`

GolangCI-Lint was created to fix the following issues with `gometalinter`:

1. Slow work: `gometalinter` usually works for minutes in average projects.
   **GolangCI-Lint works [2-7x times faster](#performance)** by [reusing work](#internals).
2. Huge memory consumption: parallel linters don't share the same program representation and can consume
   `n` times more memory (`n` - concurrency). GolangCI-Lint fixes it by sharing representation and **consumes 26% less memory**.
3. Doesn't use real bounded concurrency: if you set it to `n` it can take up to `n*n` threads because of
   forced threads in specific linters. `gometalinter` can't do anything about it because it runs linters as
   black boxes in forked processes. In GolangCI-Lint we run all linters in one process and completely control
   them. Configured concurrency will be correctly bounded.
   This issue is important because you often want to set concurrency to the CPUs count minus one to
   ensure you **do not freeze your PC** and be able to work on it while analyzing code.
4. Lack of nice output. We like how the `gcc` and `clang` compilers format their warnings: **using colors,
   printing warning lines and showing the position in line**.
5. Too many issues. GolangCI-Lint cuts a lot of issues by using default exclude list of common false-positives.
   By default, it has enabled **smart issues processing**: merge multiple issues for one line, merge issues with the
   same text or from the same linter. All of these smart processors can be configured by the user.
6. Integration into large codebases. A good way to start using linters in a large project is not to fix a plethora
   of existing issues, but to set up CI and **fix only issues in new commits**. You can use `revgrep` for it, but it's
   yet another utility to install and configure. With `golangci-lint` it's much easier: `revgrep` is already built into
   `golangci-lint` and you can use it with one option (`-n, --new` or `--new-from-rev`).
7. Installation. With `gometalinter`, you need to run a linters installation step. It's easy to forget this step and
   end up with stale linters. It also complicates CI setup. GolangCI-Lint requires **no installation of linters**.
8. **Yaml or toml config**. Gometalinter's JSON isn't convenient for config files.

### `golangci-lint` vs Running Linters Manually

1. It will be much slower because `golangci-lint` runs all linters in parallel and shares 50-80% of linters work.
2. It will have less control and more false-positives: some linters can't be properly configured without hacks.
3. It will take more time because of different usages and need of tracking of versions of `n` linters.

## Performance

Benchmarks were executed on MacBook Pro (Retina, 13-inch, Late 2013), 2,4 GHz Intel Core i5, 8 GB 1600 MHz DDR3.
It has 4 cores and concurrent linting as a default consuming all cores.
Benchmark was run (and measured) automatically, see the code
[here](https://github.com/golangci/golangci-lint/blob/master/test/bench_test.go) (`BenchmarkWithGometalinter`).

We measure peak memory usage (RSS) by tracking of processes RSS every 5 ms.

### Comparison with gometalinter

We compare golangci-lint and gometalinter in default mode, but explicitly enable all linters because of small differences in the default configuration.

```bash
$ golangci-lint run --no-config --issues-exit-code=0 --timeout=30m \
  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck
$ gometalinter --deadline=30m --vendor --cyclo-over=30 --dupl-threshold=150 \
  --exclude=<default golangci-lint excludes> --skip=testdata --skip=builtin \
  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck
  ./...
```

| Repository | GolangCI Time | GolangCI Is Faster than Gometalinter | GolangCI Memory | GolangCI eats less memory than Gometalinter |
| ---------- | ------------- | ------------------------------------ | --------------- | ------------------------------------------- |
| gometalinter repo, 4 kLoC   | 6s    | **6.4x** | 0.7GB | 33%  |
| self-repo, 4 kLoC           | 12s   | **7.5x** | 1.2GB | 41%  |
| beego, 50 kLoC              | 10s   | **4.2x** | 1.4GB | 9%   |
| hugo, 70 kLoC               | 15s   | **6.1x** | 1.6GB | 44%  |
| consul, 127 kLoC            | 58s   | **4x**   | 2.7GB | 41%  |
| terraform, 190 kLoC         | 2m13s | **1.6x** | 4.8GB | 0%   |
| go-ethereum, 250 kLoC       | 33s   | **5x**   | 3.6GB | 0%   |
| go source (`$GOROOT/src`), 1300 kLoC | 2m45s | **2x** | 4.7GB | 0% |

**On average golangci-lint is 4.6 times faster** than gometalinter. Maximum difference is in the
self-repo: **7.5 times faster**, minimum difference is in terraform source code repo: 1.8 times faster.

On average golangci-lint consumes 26% less memory.

### Why golangci-lint is faster

Golangci-lint directly calls linters (no forking) and reuses 80% of work by parsing program only once.
Read [this section](#internals) for details.

### Memory Usage of Golangci-lint

A trade-off between memory usage and execution time can be controlled by [`GOGC`](https://golang.org/pkg/runtime/#hdr-Environment_Variables) environment variable.
Less `GOGC` values trigger garbage collection more frequently and golangci-lint consumes less memory and more CPU. Below is the trade-off table for running on this repo:

|`GOGC`|Peak Memory, GB|Executon Time, s|
|------|---------------|----------------|
|`5`   |1.1            |60              |
|`10`  |1.1            |34              |
|`20`  |1.3            |25              |
|`30`  |1.6            |20.2            |
|`50`  |2.0            |17.1            |
|`80`  |2.2            |14.1            |
|`100` (default)|2.2   |13.8            |
|`off` |3.2            |9.3             |

## Internals

1. Work sharing
  The key difference with gometalinter is that golangci-lint shares work between specific linters (golint, govet, ...).
  We don't fork to call specific linter but use its API.
  For small and medium projects 50-90% of work between linters can be reused.

   * load `[]*packages.Package` by `go/packages` once

      We load program (parsing all files and type-checking) only once for all linters. For the most of linters
      it's the most heavy operation: it takes 5 seconds on 8 kLoC repo and 11 seconds on `$GOROOT/src`.
   * build `ssa.Program` once

      Some linters (megacheck, interfacer, unparam) work on SSA representation.
      Building of this representation takes 1.5 seconds on 8 kLoC repo and 6 seconds on `$GOROOT/src`.

   * parse source code and build AST once

      Parsing one source file takes 200 us on average. Parsing of all files in `$GOROOT/src` takes 2 seconds.
      Currently we parse each file more than once because it's not the bottleneck. But we already save a lot of
      extra parsing. We're planning to parse each file only once.

   * walk files and directories once

     It takes 300-1000 ms for `$GOROOT/src`.
2. Smart linters scheduling

   We schedule linters by a special algorithm which takes estimated execution time into account. It allows
   to save 10-30% of time when one of heavy linters (megacheck etc) is enabled.

3. Don't fork to run shell commands

All linters are vendored in the `/vendor` folder: their version is fixed, they are builtin
and you don't need to install them separately.

## Supported Linters

To see a list of supported linters and which linters are enabled/disabled:

```bash
golangci-lint help linters
```

### Enabled By Default Linters

- [govet](https://golang.org/cmd/vet/) - Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
- [errcheck](https://github.com/kisielk/errcheck) - Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases
- [staticcheck](https://staticcheck.io/) - Staticcheck is a go vet on steroids, applying a ton of static analysis checks
- [unused](https://github.com/dominikh/go-tools/tree/master/unused) - Checks Go code for unused constants, variables, functions and types
- [gosimple](https://github.com/dominikh/go-tools/tree/master/simple) - Linter for Go source code that specializes in simplifying a code
- [structcheck](https://github.com/opennota/check) - Finds unused struct fields
- [varcheck](https://github.com/opennota/check) - Finds unused global variables and constants
- [ineffassign](https://github.com/gordonklaus/ineffassign) - Detects when assignments to existing variables are not used
- [deadcode](https://github.com/remyoudompheng/go-misc/tree/master/deadcode) - Finds unused code
- typecheck - Like the front-end of a Go compiler, parses and type-checks Go code

### Disabled By Default Linters (`-E/--enable`)

- [bodyclose](https://github.com/timakin/bodyclose) - checks whether HTTP response body is closed successfully
- [golint](https://github.com/golang/lint) - Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
- [stylecheck](https://github.com/dominikh/go-tools/tree/master/stylecheck) - Stylecheck is a replacement for golint
- [gosec](https://github.com/securego/gosec) - Inspects source code for security problems
- [interfacer](https://github.com/mvdan/interfacer) - Linter that suggests narrower interface types
- [unconvert](https://github.com/mdempsky/unconvert) - Remove unnecessary type conversions
- [dupl](https://github.com/mibk/dupl) - Tool for code clone detection
- [goconst](https://github.com/jgautheron/goconst) - Finds repeated strings that could be replaced by a constant
- [gocyclo](https://github.com/alecthomas/gocyclo) - Computes and checks the cyclomatic complexity of functions
- [gocognit](https://github.com/uudashr/gocognit) - Computes and checks the cognitive complexity of functions
- [gofmt](https://golang.org/cmd/gofmt/) - Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) - Goimports does everything that gofmt does. Additionally it checks unused imports
- [maligned](https://github.com/mdempsky/maligned) - Tool to detect Go structs that would take less memory if their fields were sorted
- [depguard](https://github.com/OpenPeeDeeP/depguard) - Go linter that checks if package imports are in a list of acceptable packages
- [misspell](https://github.com/client9/misspell) - Finds commonly misspelled English words in comments
- [lll](https://github.com/walle/lll) - Reports long lines
- [unparam](https://github.com/mvdan/unparam) - Reports unused function parameters
- [dogsled](https://github.com/alexkohler/dogsled) - Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f())
- [nakedret](https://github.com/alexkohler/nakedret) - Finds naked returns in functions greater than a specified function length
- [prealloc](https://github.com/alexkohler/prealloc) - Finds slice declarations that could potentially be preallocated
- [scopelint](https://github.com/kyoh86/scopelint) - Scopelint checks for unpinned variables in go programs
- [gocritic](https://github.com/go-critic/go-critic) - The most opinionated Go source code linter
- [gochecknoinits](https://github.com/leighmcculloch/gochecknoinits) - Checks that no init functions are present in Go code
- [gochecknoglobals](https://github.com/leighmcculloch/gochecknoglobals) - Checks that no globals are present in Go code
- [godox](https://github.com/matoous/godox) - Tool for detection of FIXME, TODO and other comment keywords
- [funlen](https://github.com/ultraware/funlen) - Tool for detection of long functions
- [whitespace](https://github.com/ultraware/whitespace) - Tool for detection of leading and trailing whitespace
- [wsl](https://github.com/bombsimon/wsl) - Whitespace Linter - Forces you to use empty lines!
- [gomnd](https://github.com/tommy-muehle/go-mnd) - An analyzer to detect magic numbers.

## Configuration

The config file has lower priority than command-line options. If the same bool/string/int option is provided on the command-line
and in the config file, the option from command-line will be used.
Slice options (e.g. list of enabled/disabled linters) are combined from the command-line and config file.

To see a list of enabled by your configuration linters:

```bash
golangci-lint linters
```

### Command-Line Options

```bash
golangci-lint run -h
Usage:
  golangci-lint run [flags]

Flags:
      --out-format string              Format of output: colored-line-number|line-number|json|tab|checkstyle|code-climate|junit-xml (default "colored-line-number")
      --print-issued-lines             Print lines of code with issue (default true)
      --print-linter-name              Print linter name in issue line (default true)
      --modules-download-mode string   Modules download mode. If not empty, passed as -mod=<mode> to go tools
      --issues-exit-code int           Exit code when issues were found (default 1)
      --build-tags strings             Build tags
      --timeout duration               Timeout for total work (default 1m0s)
      --tests                          Analyze tests (*_test.go) (default true)
      --print-resources-usage          Print avg and max memory usage of golangci-lint and total time
  -c, --config PATH                    Read config from file path PATH
      --no-config                      Don't read config
      --skip-dirs strings              Regexps of directories to skip
      --skip-dirs-use-default          Use or not use default excluded directories:
                                         - (^|/)vendor($|/)
                                         - (^|/)third_party($|/)
                                         - (^|/)testdata($|/)
                                         - (^|/)examples($|/)
                                         - (^|/)Godeps($|/)
                                         - (^|/)builtin($|/)
                                        (default true)
      --skip-files strings             Regexps of files to skip
  -E, --enable strings                 Enable specific linter
  -D, --disable strings                Disable specific linter
      --disable-all                    Disable all linters
  -p, --presets strings                Enable presets (bugs|complexity|format|performance|style|unused) of linters. Run 'golangci-lint linters' to see them. This option implies option --disable-all
      --fast                           Run only fast linters from enabled linters set (first run won't be fast)
  -e, --exclude strings                Exclude issue by regexp
      --exclude-use-default            Use or not use default excludes:
                                         # errcheck: Almost all programs ignore errors on these functions and in most cases it's ok
                                         - Error return value of .((os\.)?std(out|err)\..*|.*Close|.*Flush|os\.Remove(All)?|.*printf?|os\.(Un)?Setenv). is not checked
                                       
                                         # golint: Annoying issue about not having a comment. The rare codebase has such comments
                                         - (comment on exported (method|function|type|const)|should have( a package)? comment|comment should be of the form)
                                       
                                         # golint: False positive when tests are defined in package 'test'
                                         - func name will be used as test\.Test.* by other packages, and that stutters; consider calling this
                                       
                                         # govet: Common false positives
                                         - (possible misuse of unsafe.Pointer|should have signature)
                                       
                                         # staticcheck: Developers tend to write in C-style with an explicit 'break' in a 'switch', so it's ok to ignore
                                         - ineffective break statement. Did you mean to break out of the outer loop
                                       
                                         # gosec: Too many false-positives on 'unsafe' usage
                                         - Use of unsafe calls should be audited
                                       
                                         # gosec: Too many false-positives for parametrized shell calls
                                         - Subprocess launch(ed with variable|ing should be audited)
                                       
                                         # gosec: Duplicated errcheck checks
                                         - G104
                                       
                                         # gosec: Too many issues in popular repos
                                         - (Expect directory permissions to be 0750 or less|Expect file permissions to be 0600 or less)
                                       
                                         # gosec: False positive is triggered by 'src, err := ioutil.ReadFile(filename)'
                                         - Potential file inclusion via variable
                                        (default true)
      --max-issues-per-linter int      Maximum issues count per one linter. Set to 0 to disable (default 50)
      --max-same-issues int            Maximum count of issues with the same text. Set to 0 to disable (default 3)
  -n, --new                            Show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only changes in HEAD~ are analyzed.
                                       It's a super-useful option for integration of golangci-lint into existing large codebase.
                                       It's not practical to fix all existing issues at the moment of integration: much better to not allow issues in new code.
                                       For CI setups, prefer --new-from-rev=HEAD~, as --new can skip linting the current patch if any scripts generate unstaged files before golangci-lint runs.
      --new-from-rev REV               Show only new issues created after git revision REV
      --new-from-patch PATH            Show only new issues created in git patch with file path PATH
      --fix                            Fix found issues (if it's supported by the linter)
  -h, --help                           help for run

Global Flags:
      --color string              Use color when printing; can be 'always', 'auto', or 'never' (default "auto")
  -j, --concurrency int           Concurrency (default NumCPU) (default 8)
      --cpu-profile-path string   Path to CPU profile output file
      --mem-profile-path string   Path to memory profile output file
      --trace-path string         Path to trace output file
  -v, --verbose                   verbose output
      --version                   Print version

```

### Config File

GolangCI-Lint looks for config files in the following paths from the current working directory:

* `.golangci.yml`
* `.golangci.toml`
* `.golangci.json`

GolangCI-Lint also searches for config files in all directories from the directory of the first analyzed path up to the root.
To see which config file is being used and where it was sourced from run golangci-lint with `-v` option.

Config options inside the file are identical to command-line options.
You can configure specific linters' options only within the config file (not the command-line).

There is a [`.golangci.example.yml`](https://github.com/golangci/golangci-lint/blob/master/.golangci.example.yml) example
config file with all supported options, their description and default value:

```yaml
# This file contains all available configuration options
# with their default values.

# options for analysis running
run:
  # default concurrency is a available CPU number
  concurrency: 4

  # timeout for analysis, e.g. 30s, 5m, default is 1m
  timeout: 1m

  # exit code when at least one issue was found, default is 1
  issues-exit-code: 1

  # include test files or not, default is true
  tests: true

  # list of build tags, all linters use it. Default is empty list.
  build-tags:
    - mytag

  # which dirs to skip: issues from them won't be reported;
  # can use regexp here: generated.*, regexp is applied on full path;
  # default value is empty list, but default dirs are skipped independently
  # from this option's value (see skip-dirs-use-default).
  skip-dirs:
    - src/external_libs
    - autogenerated_by_my_lib

  # default is true. Enables skipping of directories:
  #   vendor$, third_party$, testdata$, examples$, Godeps$, builtin$
  skip-dirs-use-default: true

  # which files to skip: they will be analyzed, but issues from them
  # won't be reported. Default value is empty list, but there is
  # no need to include all autogenerated files, we confidently recognize
  # autogenerated files. If it's not please let us know.
  skip-files:
    - ".*\\.my\\.go$"
    - lib/bad.go

  # by default isn't set. If set we pass it to "go list -mod={option}". From "go help modules":
  # If invoked with -mod=readonly, the go command is disallowed from the implicit
  # automatic updating of go.mod described above. Instead, it fails when any changes
  # to go.mod are needed. This setting is most useful to check that go.mod does
  # not need updates, such as in a continuous integration and testing system.
  # If invoked with -mod=vendor, the go command assumes that the vendor
  # directory holds the correct copies of dependencies and ignores
  # the dependency descriptions in go.mod.
  modules-download-mode: readonly|release|vendor


# output configuration options
output:
  # colored-line-number|line-number|json|tab|checkstyle|code-climate, default is "colored-line-number"
  format: colored-line-number

  # print lines of code with issue, default is true
  print-issued-lines: true

  # print linter name in the end of issue text, default is true
  print-linter-name: true


# all available settings of specific linters
linters-settings:
  errcheck:
    # report about not checking of errors in type assetions: `a := b.(MyStruct)`;
    # default is false: such cases aren't reported by default.
    check-type-assertions: false

    # report about assignment of errors to blank identifier: `num, _ := strconv.Atoi(numStr)`;
    # default is false: such cases aren't reported by default.
    check-blank: false

    # [deprecated] comma-separated list of pairs of the form pkg:regex
    # the regex is used to ignore names within pkg. (default "fmt:.*").
    # see https://github.com/kisielk/errcheck#the-deprecated-method for details
    ignore: fmt:.*,io/ioutil:^Read.*

    # path to a file containing a list of functions to exclude from checking
    # see https://github.com/kisielk/errcheck#excluding-functions for details
    exclude: /path/to/file.txt

  funlen:
    lines: 60
    statements: 40

  govet:
    # report about shadowed variables
    check-shadowing: true

    # settings per analyzer
    settings:
      printf: # analyzer name, run `go tool vet help` to see all analyzers
        funcs: # run `go tool vet help printf` to see available settings for `printf` analyzer
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf

    # enable or disable analyzers by name
    enable:
      - atomicalign
    enable-all: false
    disable:
      - shadow
    disable-all: false
  golint:
    # minimal confidence for issues, default is 0.8
    min-confidence: 0.8
  gofmt:
    # simplify code: gofmt with `-s` option, true by default
    simplify: true
  goimports:
    # put imports beginning with prefix after 3rd-party packages;
    # it's a comma-separated list of prefixes
    local-prefixes: github.com/org/project
  gocyclo:
    # minimal code complexity to report, 30 by default (but we recommend 10-20)
    min-complexity: 10
  gocognit:
    # minimal code complexity to report, 30 by default (but we recommend 10-20)
    min-complexity: 10
  maligned:
    # print struct with more effective memory layout or not, false by default
    suggest-new: true
  dupl:
    # tokens count to trigger issue, 150 by default
    threshold: 100
  goconst:
    # minimal length of string constant, 3 by default
    min-len: 3
    # minimal occurrences count to trigger, 3 by default
    min-occurrences: 3
  depguard:
    list-type: blacklist
    include-go-root: false
    packages:
      - github.com/sirupsen/logrus
    packages-with-error-message:
      # specify an error message to output when a blacklisted package is used
      - github.com/sirupsen/logrus: "logging is allowed only by logutils.Log"
  misspell:
    # Correct spellings using locale preferences for US or UK.
    # Default is to use a neutral variety of English.
    # Setting locale to US will correct the British spelling of 'colour' to 'color'.
    locale: US
    ignore-words:
      - someword
  lll:
    # max line length, lines longer will be reported. Default is 120.
    # '\t' is counted as 1 character by default, and can be changed with the tab-width option
    line-length: 120
    # tab width in spaces. Default to 1.
    tab-width: 1
  unused:
    # treat code as a program (not a library) and report unused exported identifiers; default is false.
    # XXX: if you enable this setting, unused will report a lot of false-positives in text editors:
    # if it's called for subdir of a project it can't find funcs usages. All text editor integrations
    # with golangci-lint call it on a directory with the changed file.
    check-exported: false
  unparam:
    # Inspect exported functions, default is false. Set to true if no external program/library imports your code.
    # XXX: if you enable this setting, unparam will report a lot of false-positives in text editors:
    # if it's called for subdir of a project it can't find external interfaces. All text editor integrations
    # with golangci-lint call it on a directory with the changed file.
    check-exported: false
  nakedret:
    # make an issue if func has more lines of code than this setting and it has naked returns; default is 30
    max-func-lines: 30
  prealloc:
    # XXX: we don't recommend using this linter before doing performance profiling.
    # For most programs usage of prealloc will be a premature optimization.

    # Report preallocation suggestions only on simple loops that have no returns/breaks/continues/gotos in them.
    # True by default.
    simple: true
    range-loops: true # Report preallocation suggestions on range loops, true by default
    for-loops: false # Report preallocation suggestions on for loops, false by default
  gocritic:
    # Which checks should be enabled; can't be combined with 'disabled-checks';
    # See https://go-critic.github.io/overview#checks-overview
    # To check which checks are enabled run `GL_DEBUG=gocritic golangci-lint run`
    # By default list of stable checks is used.
    enabled-checks:
      - rangeValCopy

    # Which checks should be disabled; can't be combined with 'enabled-checks'; default is empty
    disabled-checks:
      - regexpMust

    # Enable multiple checks by tags, run `GL_DEBUG=gocritic golangci-lint run` to see all tags and checks.
    # Empty list by default. See https://github.com/go-critic/go-critic#usage -> section "Tags".
    enabled-tags:
      - performance

    settings: # settings passed to gocritic
      captLocal: # must be valid enabled check name
        paramsOnly: true
      rangeValCopy:
        sizeThreshold: 32
  godox:
    # report any comments starting with keywords, this is useful for TODO or FIXME comments that
    # might be left in the code accidentally and should be resolved before merging
    keywords: # default keywords are TODO, BUG, and FIXME, these can be overwritten by this setting
      - NOTE
      - OPTIMIZE # marks code that should be optimized before merging
      - HACK # marks hack-arounds that should be removed before merging
  dogsled:
    # checks assignments with too many blank identifiers; default is 2
    max-blank-identifiers: 2

  whitespace:
    multi-if: false   # Enforces newlines (or comments) after every multi-line if statement
    multi-func: false # Enforces newlines (or comments) after every multi-line function signature
  wsl:
    # If true append is only allowed to be cuddled if appending value is
    # matching variables, fields or types on line above. Default is true.
    strict-append: true
    # Allow calls and assignments to be cuddled as long as the lines have any
    # matching variables, fields or types. Default is true.
    allow-assign-and-call: true
    # Allow multiline assignments to be cuddled. Default is true.
    allow-multiline-assign: true
    # Allow declarations (var) to be cuddled.
    allow-cuddle-declarations: false
    # Allow trailing comments in ending of blocks
    allow-trailing-comment: false
    # Force newlines in end of case at this limit (0 = never).
    force-case-trailing-whitespace: 0

  # The custom section can be used to define linter plugins to be loaded at runtime. See README doc
  #  for more info.
  custom:
    # Each custom linter should have a unique name.
     example:
      # The path to the plugin *.so. Can be absolute or local. Required for each custom linter
      path: /path/to/example.so
      # The description of the linter. Optional, just for documentation purposes.
      description: This is an example usage of a plugin linter.
      # Intended to point to the repo location of the linter. Optional, just for documentation purposes.
      original-url: github.com/golangci/example-linter

linters:
  enable:
    - megacheck
    - govet
  disable:
    - maligned
    - prealloc
  disable-all: false
  presets:
    - bugs
    - unused
  fast: false


issues:
  # List of regexps of issue texts to exclude, empty list by default.
  # But independently from this option we use default exclude patterns,
  # it can be disabled by `exclude-use-default: false`. To list all
  # excluded by default patterns execute `golangci-lint run --help`
  exclude:
    - abcdef

  # Excluding configuration per-path, per-linter, per-text and per-source
  exclude-rules:
    # Exclude some linters from running on tests files.
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec

    # Exclude known linters from partially hard-vendored code,
    # which is impossible to exclude via "nolint" comments.
    - path: internal/hmac/
      text: "weak cryptographic primitive"
      linters:
        - gosec

    # Exclude some staticcheck messages
    - linters:
        - staticcheck
      text: "SA9003:"

    # Exclude lll issues for long lines with go:generate
    - linters:
        - lll
      source: "^//go:generate "

  # Independently from option `exclude` we use default exclude patterns,
  # it can be disabled by this option. To list all
  # excluded by default patterns execute `golangci-lint run --help`.
  # Default value for this option is true.
  exclude-use-default: false

  # Maximum issues count per one linter. Set to 0 to disable. Default is 50.
  max-issues-per-linter: 0

  # Maximum count of issues with the same text. Set to 0 to disable. Default is 3.
  max-same-issues: 0

  # Show only new issues: if there are unstaged changes or untracked files,
  # only those changes are analyzed, else only changes in HEAD~ are analyzed.
  # It's a super-useful option for integration of golangci-lint into existing
  # large codebase. It's not practical to fix all existing issues at the moment
  # of integration: much better don't allow issues in new code.
  # Default is false.
  new: false

  # Show only new issues created after git revision `REV`
  new-from-rev: REV

  # Show only new issues created in git patch with set file path.
  new-from-patch: path/to/patch/file
```

It's a [.golangci.yml](https://github.com/golangci/golangci-lint/blob/master/.golangci.yml) config file of this repo: we enable more linters
than the default and have more strict settings:

```yaml
linters-settings:
  govet:
    check-shadowing: true
    settings:
      printf:
        funcs:
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf
  golint:
    min-confidence: 0
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2
  depguard:
    list-type: blacklist
    packages:
      # logging is allowed only by logutils.Log, logrus
      # is allowed to use only in logutils package
      - github.com/sirupsen/logrus
    packages-with-error-message:
      - github.com/sirupsen/logrus: "logging is allowed only by logutils.Log"
  misspell:
    locale: US
  lll:
    line-length: 140
  goimports:
    local-prefixes: github.com/golangci/golangci-lint
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
    disabled-checks:
      - dupImport # https://github.com/go-critic/go-critic/issues/845
      - ifElseChain
      - octalLiteral
      - whyNoLint
      - wrapperFunc
  funlen:
    lines: 100
    statements: 50

linters:
  # please, do not use `enable-all`: it's deprecated and will be removed soon.
  # inverted configuration with `enable-all` and `disable` is not scalable during updates of golangci-lint
  disable-all: true
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

  # don't enable:
  # - gochecknoglobals
  # - gocognit
  # - godox
  # - maligned
  # - prealloc

run:
  skip-dirs:
    - test/testdata_etc
    - internal/cache
    - internal/renameio
    - internal/robustio

# golangci.com configuration
# https://github.com/golangci/golangci/wiki/Configuration
service:
  golangci-lint-version: 1.20.x # use the fixed version to not introduce new linters unexpectedly
  prepare:
    - echo "here I can run custom commands, but no preparation needed for this repo"
```

## Custom Linters
Some people and organizations may choose to have custom made linters run as a part of golangci-lint. That functionality 
is supported through go's plugin library.

### Create a Copy of `golangci-lint` that Can Run with Plugins
In order to use plugins, you'll need a golangci-lint executable that can run them. The normal version of this project 
is built with the vendors option, which breaks plugins that have overlapping dependencies.

1. Download [golangci-lint](https://github.com/golangci/golangci-lint) source code
2. From the projects root directory, run `make vendor_free_build`
3. Copy the `golangci-lint` executable that was created to your path, project, or other location

### Configure Your Project for Linting
If you already have a linter plugin available, you can follow these steps to define it's usage in a projects 
`.golangci.yml` file. An example linter can be found at [here](https://github.com/golangci/example-plugin-linter). If you're looking for 
instructions on how to configure your own custom linter, they can be found further down.

1. If the project you want to lint does not have one already, copy the [.golangci.yml](https://github.com/golangci/golangci-lint/blob/master/.golangci.yml) to the root directory.
2. Adjust the yaml to appropriate `linters-settings:custom` entries as so:
```
linters-settings:
 custom:
  example:
   path: /example.so
   description: The description of the linter
   original-url: github.com/golangci/example-linter
```

That is all the configuration that is required to run a custom linter in your project. Custom linters are enabled by default,
but abide by the same rules as other linters. If the disable all option is specified either on command line or in 
`.golang.yml` files `linters:disable-all: true`, custom linters will be disabled; they can be re-enabled by adding them 
to the `linters:enable` list, or providing the enabled option on the command line, `golangci-lint run -Eexample`.

### To Create Your Own Custom Linter

Your linter must implement one or more `golang.org/x/tools/go/analysis.Analyzer` structs.
Your project should also use `go.mod`. All versions of libraries that overlap `golangci-lint` (including replaced 
libraries) MUST be set to the same version as `golangci-lint`. You can see the versions by running `go version -m golangci-lint`.

You'll also need to create a go file like `plugin/example.go`. This MUST be in the package `main`, and define a 
variable of name `AnalyzerPlugin`. The `AnalyzerPlugin` instance MUST implement the following interface:
```
type AnalyzerPlugin interface {
    GetAnalyzers() []*analysis.Analyzer
}
```
The type of `AnalyzerPlugin` is not important, but is by convention `type analyzerPlugin struct {}`. See 
[plugin/example.go](https://github.com/golangci/example-plugin-linter/plugin/example.go) for more info.

To build the plugin, from the root project directory, run `go build -buildmode=plugin plugin/example.go`. This will create a plugin `*.so`
file that can be copied into your project or another well known location for usage in golangci-lint.

## False Positives

False positives are inevitable, but we did our best to reduce their count. For example, we have a default enabled set of [exclude patterns](#command-line-options). If a false positive occurred you have the following choices:

1. Exclude issue by text using command-line option `-e` or config option `issues.exclude`. It's helpful when you decided to ignore all issues of this type. Also, you can use `issues.exclude-rules` config option for per-path or per-linter configuration.
2. Exclude this one issue by using special comment `//nolint` (see [the section](#nolint) below).
3. Exclude issues in path by `run.skip-dirs`, `run.skip-files` or `issues.exclude-rules` config options.

Please create [GitHub Issues here](https://github.com/golangci/golangci-lint/issues/new) if you find any false positives. We will add it to the default exclude list if it's common or we will fix underlying linter.

### Nolint

To exclude issues from all linters use `//nolint`. For example, if it's used inline (not from the beginning of the line) it excludes issues only for this line.

```go
var bad_name int //nolint
```

To exclude issues from specific linters only:

```go
var bad_name int //nolint:golint,unused
```

To exclude issues for the block of code use this directive on the beginning of a line:

```go
//nolint
func allIssuesInThisFunctionAreExcluded() *string {
  // ...
}

//nolint:govet
var (
  a int
  b int
)
```

Also, you can exclude all issues in a file by:

```go
//nolint:unparam
package pkg
```

You may add a comment explaining or justifying why `//nolint` is being used on the same line as the flag itself:

```go
//nolint:gocyclo // This legacy function is complex but the team too busy to simplify it
func someLegacyFunction() *string {
  // ...
}
```

You can see more examples of using `//nolint` in [our tests](https://github.com/golangci/golangci-lint/tree/master/pkg/result/processors/testdata) for it.

Use `//nolint` instead of `// nolint` because machine-readable comments should have no space by Go convention.

## FAQ

**How do you add a custom linter?**

You can integrate it yourself, see this [wiki page](https://github.com/golangci/golangci-lint/wiki/How-to-add-a-custom-linter) with documentation. Or you can create a [GitHub Issue](https://github.com/golangci/golangci-lint/issues/new) and we will integrate when time permits.

**It's cool to use `golangci-lint` when starting a project, but what about existing projects with large codebase? It will take days to fix all found issues**

We are sure that every project can easily integrate `golangci-lint`, even the large one. The idea is to not fix all existing issues. Fix only newly added issue: issues in new code. To do this setup CI (or better use [GolangCI](https://golangci.com)) to run `golangci-lint` with option `--new-from-rev=HEAD~1`. Also, take a look at option `--new`, but consider that CI scripts that generate unstaged files will make `--new` only point out issues in those files and not in the last commit. In that regard `--new-from-rev=HEAD~1` is safer.
By doing this you won't create new issues in your code and can choose fix existing issues (or not).

**How to use `golangci-lint` in CI (Continuous Integration)?**

You have 2 choices:

1. Use [GolangCI](https://golangci.com): this service is highly integrated with GitHub (issues are commented in the pull request) and uses a `golangci-lint` tool. For configuration use `.golangci.yml` (or toml/json).
2. Use custom CI: just run `golangci-lint` in CI and check the exit code. If it's non-zero - fail the build. The main disadvantage is that you can't see issues in pull request code and would need to view the build log, then open the referenced source file to see the context.

We don't recommend vendoring `golangci-lint` in your repo: you will get troubles updating `golangci-lint`. Please, use recommended way to install with the shell script: it's very fast.

**Do I need to run `go install`?**

No, you don't need to do it anymore.

**Which go versions are supported**
Short answer: go 1.12 and newer are officially supported.

Long answer:

1. go < 1.9 isn't supported
2. go1.9 is officially supported by golangci-lint <= v1.10.2
3. go1.10 is officially supported by golangci-lint <= 1.15.0.
4. go1.11 is officially supported by golangci-lint <= 1.17.1.
5. go1.12+ are officially supported by the latest version of golangci-lint (>= 1.18.0).

**`golangci-lint` doesn't work**

1. Please, ensure you are using the latest binary release.
2. Run it with `-v` option and check the output.
3. If it doesn't help create a [GitHub issue](https://github.com/golangci/golangci-lint/issues/new) with the output from the error and #2 above.

**Why running with `--fast` is slow on the first run?**
Because the first run caches type information. All subsequent runs will be fast.
Usually this options is used during development on local machine and compilation was already performed.

## Thanks

Thanks to all [contributors](https://github.com/golangci/golangci-lint/graphs/contributors)!
Thanks to [alecthomas/gometalinter](https://github.com/alecthomas/gometalinter) for inspiration and amazing work.
Thanks to [bradleyfalzon/revgrep](https://github.com/bradleyfalzon/revgrep) for cool diff tool.

Thanks to developers and authors of used linters:
- [timakin](https://github.com/timakin)
- [kisielk](https://github.com/kisielk)
- [golang](https://github.com/golang)
- [dominikh](https://github.com/dominikh)
- [securego](https://github.com/securego)
- [opennota](https://github.com/opennota)
- [mvdan](https://github.com/mvdan)
- [mdempsky](https://github.com/mdempsky)
- [gordonklaus](https://github.com/gordonklaus)
- [mibk](https://github.com/mibk)
- [jgautheron](https://github.com/jgautheron)
- [remyoudompheng](https://github.com/remyoudompheng)
- [alecthomas](https://github.com/alecthomas)
- [uudashr](https://github.com/uudashr)
- [OpenPeeDeeP](https://github.com/OpenPeeDeeP)
- [client9](https://github.com/client9)
- [walle](https://github.com/walle)
- [alexkohler](https://github.com/alexkohler)
- [kyoh86](https://github.com/kyoh86)
- [go-critic](https://github.com/go-critic)
- [leighmcculloch](https://github.com/leighmcculloch)
- [matoous](https://github.com/matoous)
- [ultraware](https://github.com/ultraware)
- [bombsimon](https://github.com/bombsimon)
- [tommy-muehle](https://github.com/tommy-muehle)

## Changelog

Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).
There is the most valuable changes log:

### September 2019

1. Support go1.13
2. Add new linters: `funlen`, `whitespace` (with auto-fix) and `godox`
3. Update linters: `gochecknoglobals`, `scopelint`, `gosec`
4. Provide pre-built binary for ARM and FreeBSD
5. 2. Fix false-positives in `unused`
6. Support `--skip-dirs-use-default`
7. Add support for bash completions

### July 2019

1. Fix parallel writes race condition
2. Update bodyclose with fixed panic

### June 2019

1. Treat Go source files as a plain text by `misspell`: it allows detecting issues in strings, variable names, etc.
2. Implement richer and more stable auto-fix of `misspell` issues.

### May 2019

1. Add [bodyclose](https://github.com/timakin/bodyclose) linter.
2. Support junit-xml output.

### April 2019

1. Update go-critic, new checkers were added: badCall, dupImports, evalOrder, newDeref
2. Fix staticcheck panic on packages that do not compile
3. Make install script work on Windows
4. Fix compatibility with the latest x/tools version and update golang.org/x/tools
5. Correct import path of module sourcegraph/go-diff
6. Fix `max-issues-per-linter` name
7. Fix linting of preprocessed files (e.g. `*.qtpl.go`, goyacc)
8. Enable auto-fixing when running via pre-commit

### March 2019

1. Support the newest `go vet` (with `go/analysis`)
2. Support configuration of `go vet`: e.g. you can set print functions by `linters-settings.govet.settings.printf.funcs`
3. Update megacheck (staticcheck) to 2019.1.1
4. Add [information](https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint) about controlling space-time trade-off into README
5. Exclude issues by source code line regexp by `issues.exclude-rules[i].source`
6. Build and test on go 1.12
7. Support `--color` option
8. Update x/tools to fix c++ issues
9. Include support for log level
10. Sort linters list in help commands

## Debug

You can see a verbose output of linter by using `-v` option.

If you would like to see more detailed logs you can set environment variable `GL_DEBUG` to debug `golangci-lint`.
It's value is a list of debug tags. For example, `GL_DEBUG=loader,gocritic golangci-lint run`.
Existing debug tags:

1. `gocritic` - debug `go-critic` linter;
2. `env` - debug `go env` command;
3. `loader` - debug packages loading (including `go/packages` internal debugging);
4. `autogen_exclude` - debug a filter excluding autogenerated source code;
5. `nolint` - debug a filter excluding issues by `//nolint` comments.

## Future Plans

1. Upstream all changes of forked linters.
2. Make it easy to write own linter/checker: it should take a minimum code, have perfect documentation, debugging and testing tooling.
3. Speed up SSA loading: on-disk cache and existing code profiling-optimizing.
4. Analyze (don't only filter) only new code: analyze only changed files and dependencies, make incremental analysis, caches.
5. Smart new issues detector: don't print existing issues on changed lines.
6. Minimize false-positives by fixing linters and improving testing tooling.
7. Automatic issues fixing (code rewrite, refactoring) where it's possible.
8. Documentation for every issue type.

## Contact Information

Slack channel: [#golangci-lint](https://slack.com/share/IS0TDK8RG/TEKQWjTZXxfK9Ta2G5HrnsMY/enQtODg0OTMxNjU0ODY2LWUyMTQ3NDc2MmNlNGU3NTNhYWE0Nzc3MjUyZjkxZWI3YjI5ODMwNDA1NTU3MmM2Yzc5ZjQyYTFkNThlODllN2Y).

You can contact the [author](https://github.com/jirfag) of GolangCI-Lint
by [denis@golangci.com](mailto:denis@golangci.com). Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).

## License Scan

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint?ref=badge_large)
