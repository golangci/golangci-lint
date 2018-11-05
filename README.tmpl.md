# GolangCI-Lint

[![Build Status](https://travis-ci.com/golangci/golangci-lint.svg?branch=master)](https://travis-ci.com/golangci/golangci-lint)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com)

GolangCI-Lint is a linters aggregator. It's fast: on average [5 times faster](#performance) than gometalinter.
It's [easy to integrate and use](#command-line-options), has [nice output](#quick-start) and has a minimum number of false positives.

GolangCI-Lint has [integrations](#editor-integration) with VS Code, GNU Emacs, Sublime Text.

Sponsored by [GolangCI.com](https://golangci.com): SaaS service for running linters on Github pull requests. Free for Open Source.

<a href="https://golangci.com/"><img src="docs/go.png" width="250px"></a>

* [Demo](#demo)
* [Install](#install)
* [Trusted By](#trusted-by)
* [Quick Start](#quick-start)
* [Editor Integration](#editor-integration)
* [Comparison](#comparison)
* [Performance](#performance)
* [Internals](#internals)
* [Supported Linters](#supported-linters)
* [Configuration](#configuration)
* [False Positives](#false-positives)
* [FAQ](#faq)
* [Thanks](#thanks)
* [Changelog](#changelog)
* [Future Plans](#future-plans)
* [Contact Information](#contact-information)

## Demo

<p align="center">
  <img src="./docs/demo.svg" width="100%">
</p>

Short 1.5 min video demo of analyzing [beego](https://github.com/astaxie/beego).
[![asciicast](https://asciinema.org/a/183662.png)](https://asciinema.org/a/183662)

## Install

### CI Installation

Most installations are done for CI (travis, circleci etc). It's important to have reproducible CI:
don't start to fail all builds at the same time. With golangci-lint this can happen if you
use `--enable-all` and a new linter is added or even without `--enable-all`: when one upstream linter
is upgraded.

It's highly recommended to install a fixed version of golangci-lint.
Releases are available on the [releases page](https://github.com/golangci/golangci-lint/releases).

The recommended way to install golangci-lint (replace `vX.Y.Z` with the latest
version from the [releases page](https://github.com/golangci/golangci-lint/releases)):

```bash
# binary will be $GOPATH/bin/golangci-lint
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin vX.Y.Z

# or install it into ./bin/
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z

golangci-lint --version
```

As a fallback you can also use `raw.githubusercontent.com`

```bash
# binary will be $GOPATH/bin/golangci-lint
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $GOPATH/bin vX.Y.Z

# or install it into ./bin/
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s vX.Y.Z

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s vX.Y.Z

golangci-lint --version
```

Periodically update version of golangci-lint: the project is under active development
and is constantly being improved. But please always check for newly found issues and
update if needed.

### Local Installation

It's a not recommended for your CI pipeline. Only install like this for your local development environment.

```bash
go get -u github.com/golangci/golangci-lint
cd $GOPATH/src/github.com/golangci/golangci-lint/cmd/golangci-lint
go install -ldflags "-X 'main.version=$(git describe --tags)' -X 'main.commit=$(git rev-parse --short HEAD)' -X 'main.date=$(date)'"
```

You can also install it on OSX using brew:

```bash
brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint
```

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

The following great projects use golangci-lint:

* [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser)
* [tsuru/tsuru](https://github.com/tsuru/tsuru)
* [xenolf/lego](https://github.com/xenolf/lego)
* [go-swagger/go-swagger](https://github.com/go-swagger/go-swagger)
* [kubernetes-sigs/kustomize](https://github.com/kubernetes-sigs/kustomize)
* [virtual-kubelet/virtual-kubelet](https://github.com/virtual-kubelet/virtual-kubelet)
* [alecthomas/participle](https://github.com/alecthomas/participle)
* [asobti/kube-monkey](https://github.com/asobti/kube-monkey)
* [getantibody/antibody](https://github.com/getantibody/antibody)
* [banzaicloud/pipeline](https://github.com/banzaicloud/pipeline)
* [posener/complete](https://github.com/posener/complete)
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
   * vim-go open [issue](https://github.com/fatih/vim-go/issues/1841)
   * syntastic [merged pull request](https://github.com/vim-syntastic/syntastic/pull/2190) with golangci-lint support
   * ale [merged pull request](https://github.com/w0rp/ale/pull/1890) with golangci-lint support
6. Atom - the [issue](https://github.com/joefitzgerald/go-plus/issues/819) for a go-plus.

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
$ golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck
$ gometalinter --deadline=30m --vendor --cyclo-over=30 --dupl-threshold=150 \
  --exclude=<defaul golangci-lint excludes> --skip=testdata --skip=builtin \
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

## Internals

1. Work sharing
  The key difference with gometalinter is that golangci-lint shares work between specific linters (golint, govet, ...).
  We don't fork to call specific linter but use its API.
  For small and medium projects 50-90% of work between linters can be reused.

   * load `loader.Program` once

      We load program (parsing all files and type-checking) only once for all linters. For the most of linters
      it's the most heavy operation: it takes 5 seconds on 8 kLoC repo and 11 seconds on `$GOROOT/src`.
   * build `ssa.Program` once

      Some linters (megacheck, interfacer, unparam) work on SSA representation.
      Building of this representation takes 1.5 seconds on 8 kLoC repo and 6 seconds on `$GOROOT/src`.
      `SSA` representation is used from a [fork of go-tools](https://github.com/dominikh/go-tools), not the official one.

   * parse source code and build AST once

      Parsing one source file takes 200 us on average. Parsing of all files in `$GOROOT/src` takes 2 seconds.
      Currently we parse each file more than once because it's not the bottleneck. But we already save a lot of
      extra parsing. We're planning to parse each file only once.

   * walk files and directories once

     It takes 300-1000 ms for `$GOROOT/src`.
2. Smart linters scheduling
  
   We schedule linters by a special algorithm which takes estimated execution time into account. It allows
   to save 10-30% of time when one of heavy linters (megacheck etc) is enabled.

3. Improved program loading

   We smartly use setting `TypeCheckFuncBodies` in `loader.Config` to build `loader.Program`.
   If there are no linters requiring SSA enabled we can load dependencies of analyzed code much faster
   by not analyzing their functions: we analyze only file-level declarations. It makes program loading
   10-30% faster in such cases.
4. Don't fork to run shell commands

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

1. Exclude issue by text using command-line option `-e` or config option `issues.exclude`. It's helpful when you decided to ignore all issues of this type.
2. Exclude this one issue by using special comment `//nolint[:linter1,linter2,...]` on issued line.
   Comment `//nolint` disables all issues reporting on this line. Comment e.g. `//nolint:govet` disables only govet issues for this line.
   If you would like to completely exclude all issues for some function prepend this comment
   above function:

```go
//nolint
func f() {
  ...
}
```

Please create [GitHub Issues here](https://github.com/golangci/golangci-lint/issues/new) if you find any false positives. We will add it to the default exclude list if it's common or we will fix underlying linter.

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

**Does I need to run `go install`?**

No, you don't need to do it anymore.

**Which go versions are supported**
Golangci-lint versions > 1.10.2 supports Go 1.10 and 1.11.
Golangci-lint versions <= v1.10.2 supported Go 1.9, 1.10, 1.11.

**`golangci-lint` doesn't work**

1. Update it: `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`
2. Run it with `-v` option and check the output.
3. If it doesn't help create a [GitHub issue](https://github.com/golangci/golangci-lint/issues/new) with the output from the error and #2 above.

**Why running with `--fast` is slow on the first run?**
Because the first run caches type information. All subsequent runs will be fast.
Usually this options is used during development on local machine and compilation was already performed.

## Thanks

Thanks to [alecthomas/gometalinter](https://github.com/alecthomas/gometalinter) for inspiration and amazing work.
Thanks to [bradleyfalzon/revgrep](https://github.com/bradleyfalzon/revgrep) for cool diff tool.

Thanks to developers and authors of used linters:
{{.ThanksList}}

## Changelog

There is the most valuable changes log:

### October 2018

1. Update goimports formatting
2. Use go/packages
   * A lot of linters became "fast": they are enabled by --fast now and
     work in 1-2 seconds. Only unparam, interfacer and megacheck
     are "slow" linters now.

   * Average project is analyzed 20-40% faster than before if all linters are
     enabled! If we enable all linters except unparam, interfacer and
     megacheck analysis is 10-20x faster!
3. Support goimports.local-prefix option for goimports
4. Change license from AGPL to GPL
5. Update CLA

### September 2018

1. Rename GAS to gosec
2. Drop go1.9 support
3. Support installation of golangci-lint via go modules
4. Update dockerfile to use golang 1.11
5. Add support for ignore/exclude flags in errcheck

### August 2018

1. Improve lll parsing for very long lines
2. Update Depguard with a Glob support
3. Silent output by default
4. Disable GAS (gosec) by default
5. Build golangci-lint on go1.11

### July 2018

1. Add `golangci-lint linters` command
2. Fix work with symlinks

### June 2018

1. Add support of the next linters:
   * unparam
   * misspell
   * prealloc
   * nakedret
   * lll
   * depguard
2. Smart generated files detector
3. Full `//nolint` support
4. Implement `--skip-files` and `--skip-dirs` options
5. Checkstyle output format support

### May 2018

1. Support GitHub Releases
2. Installation via Homebrew and Docker

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
by [denis@golangci.com](mailto:denis@golangci.com).

## License Scan

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint?ref=badge_large)