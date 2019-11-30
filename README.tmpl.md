# GolangCI-Lint

[![Build Status](https://travis-ci.com/golangci/golangci-lint.svg?branch=master)](https://travis-ci.com/golangci/golangci-lint)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)

GolangCI-Lint is a linters aggregator. It's fast: on average [5 times faster](#performance) than gometalinter.
It's [easy to integrate and use](#command-line-options), has [nice output](#quick-start) and has a minimum number of false positives. It supports go modules.

GolangCI-Lint has [integrations](#editor-integration) with VS Code, GNU Emacs, Sublime Text.

Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).

Sponsored by [GolangCI.com](https://golangci.com): SaaS service for running linters on Github pull requests. Free for Open Source.

<a href="https://golangci.com/"><img src="docs/go.png" width="250px"></a>

- [GolangCI-Lint](#golangci-lint)
  - [Demo](#demo)
  - [Install](#install)
    - [Binary Release](#binary-release)
    - [MacOS](#macos)
    - [By Docker](#by-docker)
    - [go get](#go-get)
  - [Trusted By](#trusted-by)
  - [Quick Start](#quick-start)
  - [Editor Integration](#editor-integration)
  - [Shell Completion](#shell-completion)
    - [Mac OS X](#mac-os-x)
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

### Binary Release

Most installations are done for CI (travis, circleci etc). It's important to have reproducible CI:
don't start to fail all builds at the same time. With golangci-lint this can happen if you
use deprecated option `--enable-all` and a new linter is added or even without `--enable-all`: when one upstream linter is upgraded.

It's highly recommended to install a specific version of golangci-lint. Releases are available on the [releases page](https://github.com/golangci/golangci-lint/releases).

Latest release: [![GitHub release](https://img.shields.io/github/release/golangci/golangci-lint.svg)]((https://github.com/golangci/golangci-lint/releases/latest))

Here is the recommended way to install golangci-lint {{.LatestVersion}}:

```bash
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin {{.LatestVersion}}

# or install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s {{.LatestVersion}}

# In alpine linux (as it does not come with curl by default)
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s {{.LatestVersion}}

golangci-lint --version
```

It is advised that you periodically update version of golangci-lint as the project is under active development
and is constantly being improved. For any problems with golangci-lint, check out recent [GitHub issues](https://github.com/golangci/golangci-lint/issues) and update if needed.

### MacOS

You can also install a binary release on MacOS using [brew](https://brew.sh/):

```bash
brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint
```

### By Docker

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:{{.LatestVersion}} golangci-lint run -v
```

### go get

Please, do not install `golangci-lint` by `go get`:

1. [`go.mod`](https://github.com/golangci/golangci-lint/blob/master/go.mod) replacement directive doesn't apply. It means you will be using patched version of `golangci-lint`.
2. it's much slower than binary installation
3. it's stability depends on your Go version (e.g. on [this compiler Go <= 1.12 bug](https://github.com/golang/go/issues/29612)).
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
{{.LintersCommandOutputEnabledOnly}}
```

and the following linters are disabled by default:

```bash
$ golangci-lint help linters
...
{{.LintersCommandOutputDisabledOnly}}
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
   * Configure [File Watcher](https://www.jetbrains.com/help/go/settings-tools-file-watchers.html) with arguments `run --print-issued-lines=false $FileDir$`.
   * Predefined File Watcher will be added in [issue](https://youtrack.jetbrains.com/issue/GO-4574).
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

### Mac OS X

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

{{.EnabledByDefaultLinters}}

### Disabled By Default Linters (`-E/--enable`)

{{.DisabledByDefaultLinters}}

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
{{.RunHelpText}}
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
{{.GolangciYamlExample}}
```

It's a [.golangci.yml](https://github.com/golangci/golangci-lint/blob/master/.golangci.yml) config file of this repo: we enable more linters
than the default and have more strict settings:

```yaml
{{.GolangciYaml}}
```

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
Short answer: go 1.12 and newer are oficially supported.

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
{{.ThanksList}}

## Changelog

{{.ChangeLog}}

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

You can contact the [author](https://github.com/jirfag) of GolangCI-Lint
by [denis@golangci.com](mailto:denis@golangci.com). Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).

## License Scan

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint?ref=badge_large)
