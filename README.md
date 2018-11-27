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
Enabled by default linters:
govet: Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string [fast: true]
errcheck: Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases [fast: true]
staticcheck: Staticcheck is a go vet on steroids, applying a ton of static analysis checks [fast: false]
unused: Checks Go code for unused constants, variables, functions and types [fast: false]
gosimple: Linter for Go source code that specializes in simplifying a code [fast: false]
structcheck: Finds an unused struct fields [fast: true]
varcheck: Finds unused global variables and constants [fast: true]
ineffassign: Detects when assignments to existing variables are not used [fast: true]
deadcode: Finds unused code [fast: true]
typecheck: Like the front-end of a Go compiler, parses and type-checks Go code [fast: true]
```

and the following linters are disabled by default:

```bash
$ golangci-lint help linters
...
Disabled by default linters:
golint: Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes [fast: true]
gosec (gas): Inspects source code for security problems [fast: true]
interfacer: Linter that suggests narrower interface types [fast: false]
unconvert: Remove unnecessary type conversions [fast: true]
dupl: Tool for code clone detection [fast: true]
goconst: Finds repeated strings that could be replaced by a constant [fast: true]
gocyclo: Computes and checks the cyclomatic complexity of functions [fast: true]
gofmt: Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification [fast: true]
goimports: Goimports does everything that gofmt does. Additionally it checks unused imports [fast: true]
maligned: Tool to detect Go structs that would take less memory if their fields were sorted [fast: true]
megacheck: 3 sub-linters in one: unused, gosimple and staticcheck [fast: false]
depguard: Go linter that checks if package imports are in a list of acceptable packages [fast: true]
misspell: Finds commonly misspelled English words in comments [fast: true]
lll: Reports long lines [fast: true]
unparam: Reports unused function parameters [fast: false]
nakedret: Finds naked returns in functions greater than a specified function length [fast: true]
prealloc: Finds slice declarations that could potentially be preallocated [fast: true]
scopelint: Scopelint checks for unpinned variables in go programs [fast: true]
gocritic: The most opinionated Go source code linter [fast: true]
gochecknoinits: Checks that no init functions are present in Go code [fast: true]
gochecknoglobals: Checks that no globals are present in Go code [fast: true]
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

- [govet](https://golang.org/cmd/vet/) - Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
- [errcheck](https://github.com/kisielk/errcheck) - Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases
- [staticcheck](https://staticcheck.io/) - Staticcheck is a go vet on steroids, applying a ton of static analysis checks
- [unused](https://github.com/dominikh/go-tools/tree/master/cmd/unused) - Checks Go code for unused constants, variables, functions and types
- [gosimple](https://github.com/dominikh/go-tools/tree/master/cmd/gosimple) - Linter for Go source code that specializes in simplifying a code
- [structcheck](https://github.com/opennota/check) - Finds an unused struct fields
- [varcheck](https://github.com/opennota/check) - Finds unused global variables and constants
- [ineffassign](https://github.com/gordonklaus/ineffassign) - Detects when assignments to existing variables are not used
- [deadcode](https://github.com/remyoudompheng/go-misc/tree/master/deadcode) - Finds unused code
- typecheck - Like the front-end of a Go compiler, parses and type-checks Go code

### Disabled By Default Linters (`-E/--enable`)

- [golint](https://github.com/golang/lint) - Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
- [gosec](https://github.com/securego/gosec) - Inspects source code for security problems
- [interfacer](https://github.com/mvdan/interfacer) - Linter that suggests narrower interface types
- [unconvert](https://github.com/mdempsky/unconvert) - Remove unnecessary type conversions
- [dupl](https://github.com/mibk/dupl) - Tool for code clone detection
- [goconst](https://github.com/jgautheron/goconst) - Finds repeated strings that could be replaced by a constant
- [gocyclo](https://github.com/alecthomas/gocyclo) - Computes and checks the cyclomatic complexity of functions
- [gofmt](https://golang.org/cmd/gofmt/) - Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) - Goimports does everything that gofmt does. Additionally it checks unused imports
- [maligned](https://github.com/mdempsky/maligned) - Tool to detect Go structs that would take less memory if their fields were sorted
- [megacheck](https://github.com/dominikh/go-tools/tree/master/cmd/megacheck) - 3 sub-linters in one: unused, gosimple and staticcheck
- [depguard](https://github.com/OpenPeeDeeP/depguard) - Go linter that checks if package imports are in a list of acceptable packages
- [misspell](https://github.com/client9/misspell) - Finds commonly misspelled English words in comments
- [lll](https://github.com/walle/lll) - Reports long lines
- [unparam](https://github.com/mvdan/unparam) - Reports unused function parameters
- [nakedret](https://github.com/alexkohler/nakedret) - Finds naked returns in functions greater than a specified function length
- [prealloc](https://github.com/alexkohler/prealloc) - Finds slice declarations that could potentially be preallocated
- [scopelint](https://github.com/kyoh86/scopelint) - Scopelint checks for unpinned variables in go programs
- [gocritic](https://github.com/go-critic/go-critic) - The most opinionated Go source code linter
- [gochecknoinits](https://github.com/leighmcculloch/gochecknoinits) - Checks that no init functions are present in Go code
- [gochecknoglobals](https://github.com/leighmcculloch/gochecknoglobals) - Checks that no globals are present in Go code

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
      --out-format string           Format of output: colored-line-number|line-number|json|tab|checkstyle (default "colored-line-number")
      --print-issued-lines          Print lines of code with issue (default true)
      --print-linter-name           Print linter name in issue line (default true)
      --issues-exit-code int        Exit code when issues were found (default 1)
      --build-tags strings          Build tags
      --deadline duration           Deadline for total work (default 1m0s)
      --tests                       Analyze tests (*_test.go) (default true)
      --print-resources-usage       Print avg and max memory usage of golangci-lint and total time
  -c, --config PATH                 Read config from file path PATH
      --no-config                   Don't read config
      --skip-dirs strings           Regexps of directories to skip
      --skip-files strings          Regexps of files to skip
  -E, --enable strings              Enable specific linter
  -D, --disable strings             Disable specific linter
      --enable-all                  Enable all linters
      --disable-all                 Disable all linters
  -p, --presets strings             Enable presets (bugs|unused|format|style|complexity|performance) of linters. Run 'golangci-lint linters' to see them. This option implies option --disable-all
      --fast                        Run only fast linters from enabled linters set (first run won't be fast)
  -e, --exclude strings             Exclude issue by regexp
      --exclude-use-default         Use or not use default excludes:
                                      # errcheck: Almost all programs ignore errors on these functions and in most cases it's ok
                                      - Error return value of .((os\.)?std(out|err)\..*|.*Close|.*Flush|os\.Remove(All)?|.*printf?|os\.(Un)?Setenv). is not checked
                                    
                                      # golint: Annoying issue about not having a comment. The rare codebase has such comments
                                      - (comment on exported (method|function|type|const)|should have( a package)? comment|comment should be of the form)
                                    
                                      # golint: False positive when tests are defined in package 'test'
                                      - func name will be used as test\.Test.* by other packages, and that stutters; consider calling this
                                    
                                      # govet: Common false positives
                                      - (possible misuse of unsafe.Pointer|should have signature)
                                    
                                      # megacheck: Developers tend to write in C-style with an explicit 'break' in a 'switch', so it's ok to ignore
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
      --max-issues-per-linter int   Maximum issues count per one linter. Set to 0 to disable (default 50)
      --max-same-issues int         Maximum count of issues with the same text. Set to 0 to disable (default 3)
  -n, --new                         Show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only changes in HEAD~ are analyzed.
                                    It's a super-useful option for integration of golangci-lint into existing large codebase.
                                    It's not practical to fix all existing issues at the moment of integration: much better to not allow issues in new code.
                                    For CI setups, prefer --new-from-rev=HEAD~, as --new can skip linting the current patch if any scripts generate unstaged files before golangci-lint runs.
      --new-from-rev REV            Show only new issues created after git revision REV
      --new-from-patch PATH         Show only new issues created in git patch with file path PATH
  -h, --help                        help for run

Global Flags:
  -j, --concurrency int           Concurrency (default NumCPU) (default 8)
      --cpu-profile-path string   Path to CPU profile output file
      --mem-profile-path string   Path to memory profile output file
  -v, --verbose                   verbose output

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
  deadline: 1m

  # exit code when at least one issue was found, default is 1
  issues-exit-code: 1

  # include test files or not, default is true
  tests: true

  # list of build tags, all linters use it. Default is empty list.
  build-tags:
    - mytag

  # which dirs to skip: they won't be analyzed;
  # can use regexp here: generated.*, regexp is applied on full path;
  # default value is empty list, but next dirs are always skipped independently
  # from this option's value:
  #   	vendor$, third_party$, testdata$, examples$, Godeps$, builtin$
  skip-dirs:
    - src/external_libs
    - autogenerated_by_my_lib

  # which files to skip: they will be analyzed, but issues from them
  # won't be reported. Default value is empty list, but there is
  # no need to include all autogenerated files, we confidently recognize
  # autogenerated files. If it's not please let us know.
  skip-files:
    - ".*\\.my\\.go$"
    - lib/bad.go


# output configuration options
output:
  # colored-line-number|line-number|json|tab|checkstyle, default is "colored-line-number"
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
  govet:
    # report about shadowed variables
    check-shadowing: true
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
      - github.com/davecgh/go-spew/spew
  misspell:
    # Correct spellings using locale preferences for US or UK.
    # Default is to use a neutral variety of English.
    # Setting locale to US will correct the British spelling of 'colour' to 'color'.
    locale: US
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
    # call graph construction algorithm (cha, rta). In general, use cha for libraries,
    # and rta for programs with main packages. Default is cha.
    algo: cha

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
    # which checks should be enabled; can't be combined with 'disabled-checks';
    # default are: [appendAssign assignOp caseOrder dupArg dupBranchBody dupCase flagDeref
    # ifElseChain regexpMust singleCaseSwitch sloppyLen switchTrue typeSwitchVar underef
    # unlambda unslice rangeValCopy defaultCaseOrder];
    # all checks list: https://github.com/go-critic/checkers
    enabled-checks:
      - rangeValCopy
    # which checks should be disabled; can't be combined with 'enabled-checks'; default is empty
    disabled-checks:
      - regexpMust
    settings: # settings passed to gocritic
      captLocal: # must be valid enabled check name
        checkLocals: true
      rangeValCopy:
        sizeThreshold: 32

linters:
  enable:
    - megacheck
    - govet
  enable-all: false
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

  # Independently from option `exclude` we use default exclude patterns,
  # it can be disabled by this option. To list all
  # excluded by default patterns execute `golangci-lint run --help`.
  # Default value for this option is true.
  exclude-use-default: false

  # Maximum issues count per one linter. Set to 0 to disable. Default is 50.
  max-per-linter: 0

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
  golint:
    min-confidence: 0
  gocyclo:
    min-complexity: 10
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
  misspell:
    locale: US
  lll:
    line-length: 140
  goimports:
    local-prefixes: github.com/golangci/golangci-lint

linters:
  enable-all: true
  disable:
    - maligned
    - prealloc
    - gosec
    - gochecknoglobals
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
- [OpenPeeDeeP](https://github.com/OpenPeeDeeP)
- [client9](https://github.com/client9)
- [walle](https://github.com/walle)
- [alexkohler](https://github.com/alexkohler)
- [kyoh86](https://github.com/kyoh86)
- [go-critic](https://github.com/go-critic)
- [leighmcculloch](https://github.com/leighmcculloch)

## Changelog

There is the most valuable changes log:

### November 2018

1. Support new linters:
   * gocritic
   * scopelint
   * gochecknointis
   * gochecknoglobals
2. Update CLA

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