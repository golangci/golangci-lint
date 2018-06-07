# GolangCI-Lint
[![Build Status](https://travis-ci.com/golangci/golangci-lint.svg?branch=master)](https://travis-ci.com/golangci/golangci-lint)

GolangCI-Lint is a linters aggregator. It's fast: on average [5 times faster](#performance) than gometalinter. It's [easy to integrate and use](#issues-options), has [nice output](#quick-start) and has a minimum number of false positives.

GolangCI-Lint has [integrations](#editor-integration) with VS Code, GNU Emacs, Sublime Text.

Sponsored by [GolangCI.com](https://golangci.com): SaaS service for running linters on Github pull requests. Free for Open Source.

<a href="https://golangci.com/"><img src="docs/go.png" width="250px"></a>

   * [Demo](#demo)
   * [Install](#install)
   * [Quick Start](#quick-start)
   * [Editor Integration](#editor-integration)
   * [Comparison](#comparison)
   * [Performance](#performance)
   * [Supported Linters](#supported-linters)
   * [Configuration](#configuration)
   * [False Positives](#false-positives)
   * [Internals](#internals)
   * [FAQ](#faq)
   * [Thanks](#thanks)
   * [Future Plans](#future-plans)
   * [Contact Information](#contact-information)

# Demo
<p align="center">
  <img src="./docs/demo.svg" width="100%">
</p>

Short 1.5 min video demo of analyzing [beego](https://github.com/astaxie/beego).
[![asciicast](https://asciinema.org/a/183662.png)](https://asciinema.org/a/183662)

# Install
## CI Installation
The most installations are done for CI (travis, circleci etc). It's important to have reproducable CI:
don't start to fail all builds at one moment. With golangci-lint this can cappen if you
use `--enable-all` and new linter is added or even without `--enable-all`: when one linter
was upgraded from the upstream.

Therefore it's highly recommended to install a fixed version of golangci-lint.
Find needed version on the [releases page](https://github.com/golangci/golangci-lint/releases).

The recommended way to install golangci-lint is the next:
```bash
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s VERSION
```

Periodically update version of golangci-lint: we do active development
and deliver a lot of improvements. But do it explicitly with checking of
newly found issues.

## Local Installation
It's a not recommended for CI method. Do it only for the local development.
```bash
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
```

You can also install it by brew:
```bash
brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint
```

# Quick Start
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
Directories are analyzed NOT recursively, to analyze them recursively append `/...` to their path.

GolangCI-Lint can be used with zero configuration. By default next linters are enabled:
```
$ golangci-lint linters
Enabled by default linters:
govet: Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string [fast: true]
errcheck: Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases [fast: false]
staticcheck: Staticcheck is a go vet on steroids, applying a ton of static analysis checks [fast: false]
unused: Checks Go code for unused constants, variables, functions and types [fast: false]
gosimple: Linter for Go source code that specializes in simplifying a code [fast: false]
gas: Inspects source code for security problems [fast: false]
structcheck: Finds an unused struct fields [fast: false]
varcheck: Finds unused global variables and constants [fast: false]
ineffassign: Detects when assignments to existing variables are not used [fast: true]
deadcode: Finds unused code [fast: false]
typecheck: Like the front-end of a Go compiler, parses and type-checks Go code [fast: false]
```

and next linters are disabled by default:
```
$ golangci-lint linters
...
Disabled by default linters:
golint: Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes [fast: true]
interfacer: Linter that suggests narrower interface types [fast: false]
unconvert: Remove unnecessary type conversions [fast: false]
dupl: Tool for code clone detection [fast: true]
goconst: Finds repeated strings that could be replaced by a constant [fast: true]
gocyclo: Computes and checks the cyclomatic complexity of functions [fast: true]
gofmt: Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification [fast: true]
goimports: Goimports does everything that gofmt does. Additionally it checks unused imports [fast: true]
maligned: Tool to detect Go structs that would take less memory if their fields were sorted [fast: false]
megacheck: 3 sub-linters in one: unused, gosimple and staticcheck [fast: false]
depguard: Go linter that checks if package imports are in a list of acceptable packages [fast: false]
```

Pass `-E/--enable` to enable linter and `-D/--disable` to disable:
```bash
$ golangci-lint run --disable-all -E errcheck
```

# Editor Integration
1. [Go for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go).
2. Sublime Text - [plugin](https://github.com/alecthomas/SublimeLinter-contrib-golang-cilint) for SublimeLinter.
3. GNU Emacs - [flycheck checker](https://github.com/weijiangan/flycheck-golangci-lint).
4. Vim - [issue](https://github.com/fatih/vim-go/issues/1841) for vim-go.

# Comparison
## `golangci-lint` vs `gometalinter`
GolangCI-Lint was created to fix next issues with `gometalinter`:
1. Slow work: `gometalinter` usually works for minutes in average projects. **GolangCI-Lint works [2-7x times faster](#performance)** by [reusing work](#internals).
2. Huge memory consumption: parallel linters don't share the same program representation and can eat `n` times more memory (`n` - concurrency). GolangCI-Lint fixes it by sharing representation and **eats 1.35x less memory**.
3. Can't set honest concurrency: if you set it to `n` it can take up to `n*n` threads because of forced threads in specific linters. `gometalinter` can't do anything about it, because it runs linters as black-boxes in forked processes. In GolangCI-Lint we run all linters in one process and fully control them. Configured concurrency will be honest.
This issue is important because often you'd like to set concurrency to CPUs count minus one to **not freeze your PC** and be able to work on it while analyzing code.
4. Lack of nice output. We like how compilers `gcc` and `clang` format their warnings: **using colors, printing of warned line and showing position in line**.
5. Too many issues. GolangCI-Lint cuts a lot of issues by using default exclude list of common false-positives. Also, it has enabled by default **smart issues processing**: merge multiple issues for one line, merge issues with the same text or from the same linter. All of these smart processors can be configured by the user.
6. Integration to large codebases. A good way to start using linters in a large project is not to fix all hundreds on existing issues, but setup CI and **fix only issues in new commits**. You can use `revgrep` for it, but it's yet another utility to install and configure. With `golangci-lint` it's much easier: `revgrep` is already built into `golangci-lint` and you can use it with one option (`-n, --new` or `--new-from-rev`).
7. Installation. With `gometalinter`, you need to run linters installation step. It's easy to forget this step and have stale linters. It also complicates CI setup. GolangCI-Lint requires **no installation of linters**.
8. **Yaml or toml config**. Gometalinter's JSON isn't convenient for configuration files.

## `golangci-lint` vs Run Needed Linters Manually
1. It will be much slower because `golangci-lint` runs all linters in parallel and shares 50-80% of linters work.
2. It will have less control and more false-positives: some linters can't be properly configured without hacks.
3. It will take more time because of different usages and need of tracking of versions of `n` linters.

# Performance
Benchmarks were executed on MacBook Pro (Retina, 13-inch, Late 2013), 2,4 GHz Intel Core i5, 8 GB 1600 MHz DDR3.
It has 4 cores and concurrency for linters was default: number of cores.
Benchmark runs and measures timings automatically, it's code is
[here](https://github.com/golangci/golangci-lint/blob/master/test/bench_test.go) (`BenchmarkWithGometalinter`).

We measure peak memory usage (RSS) by tracking of processes RSS every 5 ms.

## Comparison with gometalinter
We compare golangci-lint and gometalinter in default mode, but explicitly specify all linters to enable because of small differences in the default configuration.
```bash
$ golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
	--disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
	--enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
	--enable=interfacer --enable=unconvert --enable=goconst --enable=gas --enable=megacheck
$ gometalinter --deadline=30m --vendor --cyclo-over=30 --dupl-threshold=150 \
	--exclude=<defaul golangci-lint excludes> --skip=testdata --skip=builtin \
	--disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
	--enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
	--enable=interfacer --enable=unconvert --enable=goconst --enable=gas --enable=megacheck
	./...
```

| Repository | GolangCI Time | GolangCI Is Faster than Gometalinter | GolangCI Memory | GolangCI eats less memory than Gometalinter |
| ---------- | ------------- | ------------------------------------ | --------------- | ------------------------------------------- |
| gometalinter repo, 4 kLoC   | 6s    | **6.4x** | 0.7GB | 1.5x |
| self-repo, 4 kLoC           | 12s   | **7.5x** | 1.2GB | 1.7x |
| beego, 50 kLoC              | 10s   | **4.2x** | 1.4GB | 1.1x |
| hugo, 70 kLoC               | 15s   | **6.1x** | 1.6GB | 1.8x |
| consul, 127 kLoC            | 58s   | **4x**   | 2.7GB | 1.7x |
| terraform, 190 kLoC         | 2m13s | **1.6x** | 4.8GB | 1x   |
| go-ethereum, 250 kLoC       | 33s   | **5x**   | 3.6GB | 1x   |
| go source, 1300 kLoC        | 2m45s | **2x**   | 4.7GB | 1x   |


**On average golangci-lint is 4.6 times faster** than gometalinter. Maximum difference is in the
self-repo: **7.5 times faster**, minimum difference is in terraform source code repo: 1.8 times faster.

On average golangci-lint consumes 1.35 times less memory.

# Supported Linters
To see a list of supported linters and which linters are enabled/disabled by default execute a command
```
golangci-lint linters
```

## Enabled By Default Linters
- [govet](https://golang.org/cmd/vet/) - Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
- [errcheck](https://github.com/kisielk/errcheck) - Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases
- [staticcheck](https://staticcheck.io/) - Staticcheck is a go vet on steroids, applying a ton of static analysis checks
- [unused](https://github.com/dominikh/go-tools/tree/master/cmd/unused) - Checks Go code for unused constants, variables, functions and types
- [gosimple](https://github.com/dominikh/go-tools/tree/master/cmd/gosimple) - Linter for Go source code that specializes in simplifying a code
- [gas](https://github.com/GoASTScanner/gas) - Inspects source code for security problems
- [structcheck](https://github.com/opennota/check) - Finds an unused struct fields
- [varcheck](https://github.com/opennota/check) - Finds unused global variables and constants
- [ineffassign](https://github.com/gordonklaus/ineffassign) - Detects when assignments to existing variables are not used
- [deadcode](https://github.com/remyoudompheng/go-misc/tree/master/deadcode) - Finds unused code
- typecheck - Like the front-end of a Go compiler, parses and type-checks Go code

## Disabled By Default Linters (`-E/--enable`)
- [golint](https://github.com/golang/lint) - Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
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

# Configuration
Configuration file has lower priority than command-line: if the same bool/string/int option defined in the command-line
and in the configuration file, option from command-line will be used.
Slice options (e.g. list of enabled/disabled linters) are combined from the command-line and configuration file.

## Command-Line Options
```
golangci-lint run -h
Usage:
  golangci-lint run [flags]

Flags:
      --out-format string           Format of output: colored-line-number|line-number|json|tab (default "colored-line-number")
      --print-issued-lines          Print lines of code with issue (default true)
      --print-linter-name           Print linter name in issue line (default true)
      --issues-exit-code int        Exit code when issues were found (default 1)
      --build-tags strings          Build tags (not all linters support them)
      --deadline duration           Deadline for total work (default 1m0s)
      --tests                       Analyze tests (*_test.go) (default true)
      --print-resources-usage       Print avg and max memory usage of golangci-lint and total time
  -c, --config PATH                 Read config from file path PATH
      --no-config                   Don't read config
      --skip-dirs strings           Regexps of directory names to skip
      --skip-files strings          Regexps of file names to skip
  -E, --enable strings              Enable specific linter
  -D, --disable strings             Disable specific linter
      --enable-all                  Enable all linters
      --disable-all                 Disable all linters
  -p, --presets strings             Enable presets (bugs|unused|format|style|complexity|performance) of linters. Run 'golangci-lint linters' to see them. This option implies option --disable-all
      --fast                        Run only fast linters from enabled linters set
  -e, --exclude strings             Exclude issue by regexp
      --exclude-use-default         Use or not use default excludes:
                                      # errcheck: Almost all programs ignore errors on these functions and in most cases it's ok
                                      - Error return value of .((os\.)?std(out|err)\..*|.*Close|.*Flush|os\.Remove(All)?|.*printf?|os\.(Un)?Setenv). is not checked
                                    
                                      # golint: Annoying issue about not having a comment. The rare codebase has such comments
                                      - (should have comment|comment on exported method|should have a package comment)
                                    
                                      # golint: False positive when tests are defined in package 'test'
                                      - func name will be used as test\.Test.* by other packages, and that stutters; consider calling this
                                    
                                      # gas: Too many false-positives on 'unsafe' usage
                                      - Use of unsafe calls should be audited
                                    
                                      # gas: Too many false-positives for parametrized shell calls
                                      - Subprocess launch(ed with variable|ing should be audited)
                                    
                                      # gas: Duplicated errcheck checks
                                      - G104
                                    
                                      # gas: Too many issues in popular repos
                                      - (Expect directory permissions to be 0750 or less|Expect file permissions to be 0600 or less)
                                    
                                      # gas: False positive is triggered by 'src, err := ioutil.ReadFile(filename)'
                                      - Potential file inclusion via variable
                                    
                                      # govet: Common false positives
                                      - (possible misuse of unsafe.Pointer|should have signature)
                                    
                                      # megacheck: Developers tend to write in C-style with an explicit 'break' in a 'switch', so it's ok to ignore
                                      - ineffective break statement. Did you mean to break out of the outer loop
                                     (default true)
      --max-issues-per-linter int   Maximum issues count per one linter. Set to 0 to disable (default 50)
      --max-same-issues int         Maximum count of issues with the same text. Set to 0 to disable (default 3)
  -n, --new                         Show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only changes in HEAD~ are analyzed.
                                    It's a super-useful option for integration of golangci-lint into existing large codebase.
                                    It's not practical to fix all existing issues at the moment of integration: much better don't allow issues in new code
      --new-from-rev REV            Show only new issues created after git revision REV
      --new-from-patch PATH         Show only new issues created in git patch with file path PATH
  -h, --help                        help for run

Global Flags:
  -j, --concurrency int           Concurrency (default NumCPU) (default 8)
      --cpu-profile-path string   Path to CPU profile output file
      --mem-profile-path string   Path to memory profile output file
  -v, --verbose                   verbose output

```

## Configuration File
GolangCI-Lint looks for next config paths in the current directory:
- `.golangci.yml`
- `.golangci.toml`
- `.golangci.json`

GolangCI-Lint also searches config file in all directories from directory of the first analyzed path up to the root.
To see which config file is used and where it was searched run golangci-lint with `-v` option.

Configuration options inside the file are identical to command-line options.
You can configure specific linters options only within configuration file, it can't be done with command-line.

There is a [`.golangci.yml`](https://github.com/golangci/golangci-lint/blob/master/.golangci.example.yml) with all supported options.

It's a [.golangci.yml](https://github.com/golangci/golangci-lint/blob/master/.golangci.yml) of this repo: we enable more linters
than by default and make their settings more strict:
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

linters:
  enable-all: true
  disable:
    - maligned
```

# False Positives
False positives are inevitable, but we did our best to reduce their count. For example, we have an enabled by default set of [exclude patterns](#issues-options). If false positive occurred you have next choices:
1. Exclude issue by text using command-line option `-e` or config option `issues.exclude`. It's helpful when you decided to ignore all issues of this type.
2. Exclude this one issue by using special comment `// nolint[:linter1,linter2,...]` on issued line.
Comment `// nolint` disables all issues reporting on this line. Comment e.g. `// nolint:govet` disables only govet issues for this line.
If you would like to completely exclude all issues for some function prepend this comment
above function:
```go
//nolint
func f() {
  ...
}
```

Please create [GitHub Issues here](https://github.com/golangci/golangci-lint/issues/new) about found false positives. We will add it to default exclude list if it's common or we will fix underlying linter.

# Internals
The key difference with gometalinter is that golangci-lint shares work between specific linters (golint, govet, ...).
For small and medium projects 50-80% of work between linters can be reused.
Now we share `loader.Program` and `SSA` representation building. `SSA` representation is used from
a [fork of go-tools](https://github.com/dominikh/go-tools), not the official one. Also, we are going to
reuse `AST` parsing and traversal.

We don't fork to call specific linter but use its API. We forked GitHub repos of almost all linters
to make API. It also allows us to be more performant and control actual count of used threads.

All linters are vendored in `/vendor` folder: their version is fixed, they are builtin
and you don't need to install them separately.

We use chains for issues and independent processors to post-process them: exclude issues by limits,
nolint comment, diff, regexps; prettify paths etc.

We use `cobra` for command-line action.

# FAQ
**Q: How to add custom linter?**

A: You can integrate it yourself, see this [wiki page](https://github.com/golangci/golangci-lint/wiki/How-to-add-a-custom-linter) with documentation. Or you can create [GitHub Issue](https://github.com/golangci/golangci-lint/issues/new) and we will integrate it soon.

**Q: It's cool to use `golangci-lint` when starting a project, but what about existing projects with large codebase? It will take days to fix all found issues**

A: We are sure that every project can easily integrate `golangci-lint`, even the large one. The idea is to not fix all existing issues. Fix only newly added issue: issues in new code. To do this setup CI (or better use [GolangCI](https://golangci.com) to run `golangci-lint` with option `--new-from-rev=origin/master`. Also, take a look at option `-n`.
By doing this you won't create new issues in code and can smoothly fix existing issues (or not).

**Q: How to use `golangci-lint` in CI (Continuous Integration)?**

A: You have 2 choices:
1. Use [GolangCI](https://golangci.com): this service is highly integrated with GitHub (issues are commented in the pull request) and uses a `golangci-lint` tool. For configuration use `.golangci.yml` (or toml/json).
2. Use custom CI: just run `golangci-lint` in CI and check the exit code. If it's non-zero - fail the build. The main disadvantage is that you can't see found issues in pull request code and should view build log, then open needed source file to see a context.
If you'd like to vendor `golangci-lint` in your repo, run:
```bash
go get -u github.com/golang/dep/cmd/dep
dep init
dep ensure -v -add github.com/golangci/golangci-lint/cmd/golangci-lint
```
Then add these lines to your `Gopkg.toml` file, so `dep ensure -update` won't delete the vendored `golangci-lint` code.
```toml
required = [
  "github.com/golangci/golangci-lint/cmd/golangci-lint",
]
```
In your CI scripts, install the vendored `golangci-lint` like this:
```bash
go install ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/`
```
Vendoring `golangci-lint` saves a network request, potentially making your CI system a little more reliable.

**Q: `golangci-lint` doesn't work**
1. Update it: `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`
2. Run it with `-v` option and check the output.
3. If it doesn't help create [GitHub issue](https://github.com/golangci/golangci-lint/issues/new) with the output.

# Thanks
Thanks to [alecthomas/gometalinter](https://github.com/alecthomas/gometalinter) for inspiration and amazing work.
Thanks to [bradleyfalzon/revgrep](https://github.com/bradleyfalzon/revgrep) for cool diff tool.

Thanks to developers and authors of used linters:
- [kisielk](https://github.com/kisielk)
- [golang](https://github.com/golang)
- [dominikh](https://github.com/dominikh)
- [GoASTScanner](https://github.com/GoASTScanner)
- [opennota](https://github.com/opennota)
- [mvdan](https://github.com/mvdan)
- [mdempsky](https://github.com/mdempsky)
- [gordonklaus](https://github.com/gordonklaus)
- [mibk](https://github.com/mibk)
- [jgautheron](https://github.com/jgautheron)
- [remyoudompheng](https://github.com/remyoudompheng)
- [alecthomas](https://github.com/alecthomas)
- [OpenPeeDeeP](https://github.com/OpenPeeDeeP)

# Future Plans
1. Upstream all changes of forked linters.
2. Fully integrate all used linters: make a common interface and reuse 100% of what can be reused: AST traversal, packages preparation etc.
3. Make it easy to write own linter/checker: it should take a minimum code, have perfect documentation, debugging and testing tooling.
4. Speedup packages loading (dig into [loader](golang.org/x/tools/go/loader)): on-disk cache and existing code profiling-optimizing.
5. Analyze (don't only filter) only new code: analyze only changed files and dependencies, make incremental analysis, caches.
6. Smart new issues detector: don't print existing issues on changed lines.
7. Integration with Text Editors. On-the-fly code analysis for text editors: it should be super-fast.
8. Minimize false-positives by fixing linters and improving testing tooling.
9. Automatic issues fixing (code rewrite, refactoring) where it's possible.
10. Documentation for every issue type.

# Contact Information
You can contact the [author](https://github.com/jirfag) of GolangCI-Lint
by [denis@golangci.com](mailto:denis@golangci.com).
