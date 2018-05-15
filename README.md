# GolangCI-Lint
[![Build Status](https://travis-ci.com/golangci/golangci-lint.svg?branch=master)](https://travis-ci.com/golangci/golangci-lint)

GolangCI-Lint is a linters aggregator. It is fast, easy to integrate and use, has nice output and has minimum count of false positives.

Sponsored by [GolangCI.com](https://golangci.com): SaaS service for running linters on Github pull requests. Free for Open Source.

<a href="https://golangci.com/"><img src="docs/go.png" width="250px"></a>

// TOC Here at the end

# Install
```bash
go get -u gopkg.in/golangci/golangci-lint.v1/cmd/golangci-lint
```

# Quick Start
To run golangci-lint execute:
```bash
golangci-lint run
```
![Screenshot of sample output](docs/run_screenshot.png)

It's and equivalent of executing:
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
govet: Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
errcheck: Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases
staticcheck: Staticcheck is go vet on steroids, applying a ton of static analysis checks
unused: Checks Go code for unused constants, variables, functions and types
gosimple: Linter for Go source code that specialises on simplifying code
gas: Inspects source code for security problems
structcheck: Finds unused struct fields
varcheck: Finds unused global variables and constants
ineffassign: Detects when assignments to existing variables are not used
deadcode: Finds unused code
```

and next linters are disabled by default:
```
$ golangci-lint linters
...
Disabled by default linters:
golint: Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
interfacer: Linter that suggests narrower interface types
unconvert: Remove unnecessary type conversions
dupl: Tool for code clone detection
goconst: Finds repeated strings that could be replaced by a constant
gocyclo: Computes and checks the cyclomatic complexity of functions
gofmt: Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification
goimports: Goimports does everything that gofmt does. Additionally it checks unused imports
maligned: Tool to detect Go structs that would take less memory if their fields were sorted
megacheck: 3 sub-linters in one: unused, gosimple and staticcheck
```

Pass `-E/--enable` to enable linter and `-D/--disable` to disable:
```bash
$ golangci-lint run --disable-all -E errcheck
```

# Comparison
## `gometalinter`
GolangCI-Lint was created to fix next issues with `gometalinter`:
1. Slow work: `gometalinter` usually works for minutes in average projects. GolangCI-Lint works [2-10x times faster](#benchmarks) by [reusing work](#internals).
2. Huge memory consumption: parallel linters don't share the same program representation and can eat `n` times more memory (`n` - concurrency). GolangCI-Lint fixes it by sharing representation.
3. Can't set honest concurrency: if you set it to `n` it can take `n+x` threads because of forced threads in specific linters. `gometalinter` can't do anything about it, because it runs linters as black-boxes in forked processes. In GolangCI-Lint we run all linters in one process and fully control them. Configured concurrency will be honest.
This issue is important because often you'd like to set concurrency to CPUs count minus one to save one CPU for example for IDE. It concurrency isn't correct you will have troubles using IDE while analyzing code.
4. Lack of nice output. We like how compilers `gcc` and `clang` format their warnings: using colors, printing of warned line and showing position in line.
5. Too many issues. GolangCI-Lint cuts a lot of issues by using default exclude list of common false-positives. Also it has enabled by default smart issues processing: merge multiple issues for one line, merge issues with the same text or from the same linter. All of these smart processors can be configured by user.
6. Installation. With `gometalinter` you need to run linters installation step. It's easy to forget this step and have stale linters.
7. Integration to large codebases. You can use `revgrep`, but it's yet another utility to install. With `golangci-lint` it's much easier: `revgrep` is already built into `golangci-lint` and you use it with one option (`-n, --new` or `--new-from-rev`).
8. Yaml or toml config. JSON isn't convenient for configuration files.

## Run Needed Linters Manually
1. It will be slower because `golangci-lint` shares 50-90% of linters work.
2. It will have fewer control and more false-positives: some linters can't be properly configured without hacks.
3. It will take more time because of different usages and need of tracking of version of `n` linters.

# Performance
## Benchmarks
```
BenchmarkWithGometalinter/self_repo/gometalinter_--fast_(4098_lines_of_code)-4                30        1482617961 ns/op
BenchmarkWithGometalinter/self_repo/golangci-lint_fast_(4098_lines_of_code)-4                100         414381899 ns/op
BenchmarkWithGometalinter/self_repo/gometalinter_(4098_lines_of_code)-4                        1        39304954722 ns/op
BenchmarkWithGometalinter/self_repo/golangci-lint_(4098_lines_of_code)-4                       5        8290405036 ns/op
```
## Internals

# Supported Linters
To see a list of supported linters and which linters are enabled/disabled by default execute command
```bash
golangci-lint linters
```

## Enabled By Default Linters
- [go vet](https://golang.org/cmd/vet/) - Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
- [errcheck](https://github.com/kisielk/errcheck): Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases
- [staticcheck](https://staticcheck.io/): Staticcheck is go vet on steroids, applying a ton of static analysis checks
- [unused](https://github.com/dominikh/go-tools/tree/master/cmd/unused): Checks Go code for unused constants, variables, functions and types
- [gosimple](https://github.com/dominikh/go-tools/tree/master/cmd/gosimple): Linter for Go source code that specialises on simplifying code
- [gas](https://github.com/GoASTScanner/gas): Inspects source code for security problems
- [structcheck](https://github.com/opennota/check): Finds unused struct fields
- [varcheck](https://github.com/opennota/check): Finds unused global variables and constants
- [ineffassign](https://github.com/gordonklaus/ineffassign): Detects when assignments to existing variables are not used
- [deadcode](https://github.com/remyoudompheng/go-misc/tree/master/deadcode): Finds unused code

## Disabled By Default Linters (`-E/--enable`)
- [golint](https://github.com/golang/lint): Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
- [interfacer](https://github.com/mvdan/interfacer): Linter that suggests narrower interface types
- [unconvert](https://github.com/mdempsky/unconvert): Remove unnecessary type conversions
- [dupl](https://github.com/mibk/dupl): Tool for code clone detection
- [goconst](https://github.com/jgautheron/goconst): Finds repeated strings that could be replaced by a constant
- [gocyclo](https://github.com/alecthomas/gocyclo): Computes and checks the cyclomatic complexity of functions
- [gofmt](https://golang.org/cmd/gofmt/): Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports): Goimports does everything that gofmt does. Additionally it checks unused imports
- [maligned](https://github.com/mdempsky/maligned): Tool to detect Go structs that would take less memory if their fields were sorted
- [megacheck](https://github.com/dominikh/go-tools/tree/master/cmd/megacheck): 3 sub-linters in one: unused, gosimple and staticcheck

# Configuration
## Command-Line Options
Run next command to see their description and defaults.
```bash
golangci-lint run -h
```

### Run Options
- `-c, --config` - path to [config file](#configuration-file) if you don't like using default config path `.golangci.(yml|toml|json)`.
- `-j, --concurrency` - number of threads used. By default it's a number of CPUs. Unlike `gometalinter`, it's an honest value, since we do not fork linter processes.
- `--build-tags` - build tags to take into account.
- `--issues-exit-code` - exit code if issues were found. Default is `1`.
- `--deadline` - timeout for running golangci-lint, `1m` by default.
- `--tests` - analyze `*_test.go` files. It's `false` by default.
- `-v, --verbose` - enable verbose output. Use this options to see which linters were enabled, to see timings of steps and another helpful information.
- `--print-resources-usage` - print memory usage and total time elapsed.

### Linters
- `-E, --enable` - enable specific linter. You can pass option multiple times or use comma:
```bash
golangci-lint run --disable-all -E golint -E govet -E errcheck
golangci-lint run --disable-all --enable golint,govet,errcheck
```
- `-D, --disable` - disable specific linter. Similar to enabling option.
- `--enable-all` - enable all supported linters.
- `--disable-all` - disable all supported linters.
- `-p, --presets` - enable specific presets. To list all presets run
```bash
$ golangci-lint linters
...
Linters presets:
bugs: govet, errcheck, staticcheck, gas, megacheck
unused: unused, structcheck, varcheck, ineffassign, deadcode, megacheck
format: gofmt, goimports
style: golint, gosimple, interfacer, unconvert, dupl, goconst, megacheck
complexity: gocyclo
performance: maligned
```
Usage example:
```bash
$ golangci-lint run -v --disable-all -p bugs,style,complexity,format
INFO[0000] Active linters: [govet goconst gocyclo gofmt gas dupl goimports megacheck interfacer unconvert errcheck golint]
```
- `--fast` - run only fast linters from enabled set of linters. To find out which linters are fast run `golangci-lint linters`.

### Linters Options
- `--errcheck.check-type-assertions` - errcheck: check for ignored type assertion results. Disabled by default.
- `--errcheck.check-blank` - errcheck: check for errors assigned to blank identifier: `_ = errFunc()`. Disabled by default
- `--govet.check-shadowing` - govet: check for shadowed variables. Disabled by default.
- `--golint.min-confidence` - golint: minimum confidence of a problem to print it. Default is `0.8`.
- `--gofmt.simplify` - gofmt: simplify code (`gofmt -s`), enabled by default.
- `--gocyclo.min-complexity` - gocyclo: minimal complexity of function to report it. Default is `30` (it's very high limit).
- `--maligned.suggest-new` - Maligned: print suggested more optimal struct fields ordering. Disabled by default. Example:
```
crypto/tls/ticket.go:20: struct of size 64 bytes could be of size 56 bytes:
struct{
	masterSecret	[]byte,
	certificates	[][]byte,
	vers        	uint16,
	cipherSuite 	uint16,
	usedOldKey  	bool,
}
```
- `--dupl.threshold` - dupl: Minimal threshold to detect copy-paste, `150` by default.
- `--goconst.min-len` - goconst: minimum constant string length, `3` by default.
- `--goconst.min-occurrences` - goconst: minimum occurences of constant string count to trigger issue. Default is `3`.

### Issues Options
- `-n, --new` - show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only changes in HEAD~ are analyzed. It's a superuseful option for integration `golangci-lint` into existing large codebase. It's not practical to fix all existing issues at the moment of integration: much better don't allow issues in new code. Disabled by default.
- `--new-from-rev` - show only new issues created after specified git revision.
- `--new-from-patch` - show only new issues created in git patch with specified file path.
- `-e, --exclude` - exclude issue by regexp on issue text.
- `--exclude-use-default` - use or not use default excludes. We tested our linter on large codebases and marked common false positives. By default we ignore common false positives by next regexps:
  - `Error return value of .((os\.)?std(out|err)\..*|.*Close|os\.Remove(All)?|.*printf?|os\.(Un)?Setenv). is not checked` - ercheck: almost all programs ignore errors on these functions and in most cases it's ok.
  - `(should have comment|comment on exported method)` - golint: annoying issues about not having a comment. Rare codebase has such comments.
  - `G103:` - gas: `Use of unsafe calls should be audited`
  - `G104:` - gas: `disable what errcheck does: it reports on Close etc`
  - `G204:` - gas: `Subprocess launching should be audited: too lot false - positives`
  - `G301:` - gas: `Expect directory permissions to be 0750 or less`
  - `G302:` - gas: `Expect file permissions to be 0600 or less`
  - `G304:` - gas: ``Potential file inclusion via variable: `src, err := ioutil.ReadFile(filename)`.``

  - `(possible misuse of unsafe.Pointer|should have signature)` - common false positives by govet.
  - `ineffective break statement. Did you mean to break out of the outer loop` - megacheck: developers tend to write in C-style with explicit `break` in switch, so it's ok to ignore.

  Use option `--exclude-use-default=false` to disable these default exclude regexps.
- `--max-issues-per-linter` - maximum issues count per one linter. Set to `0` to disable. Default value is `50` to not being annoying.
- `--max-same-issues` - maximum count of issues with the same text. Set to 0 to disable. Default value is `3` to not being annoying.

### Output Options
- `--out-format` - format of output: `colored-line-number|line-number|json`, default is `colored-line-number`.
- `--print-issued-lines` - print line of source code where issue occured. Enabled by default.
- `--print-linter-name` - print linter name in issue line. Enabled by default.
- `--print-welcome` - print welcome message. Enabled by default.

## Configuration File
GolangCI-Lint looks for next config paths in current directory:
- `.golangci.yml`
- `.golangci.toml`
- `.golangci.json`

Configuration options inside file are identical to command-line options:
- [Example `.golangci.yml`](https://github.com/golangci/golangci-lint/blob/master/.golangci.example.yml) with all supported options.
- [.golangci.yml](https://github.com/golangci/golangci-lint/blob/master/.golangci.yml) of this repo: we enable more linters than by default and make their settings more strict.

# False Positives
False positives are inevitable, but we did our best to reduce their count. If false postive occured you have next choices:
1. Exclude issue text using command-line option `-e` or config option `issues.exclude`. It's helpful when you decided to ignore all issues of this type.
2. Exclude this one issue by using special comment `// nolint[:linter1,linter2,...]` on issued line.
Comment `// nolint` disables all issues reporting on this line. Comment e.g. `// nolint:govet` disables only govet issues for this line.

Please create [GitHub Issues here](https://github.com/golangci/golangci-lint/issues/new) about found false positives. We will add it to default exclude list if it's common or we will even fix underlying linter.

# FAQ
**Q: How to add custom linter?**

A: You can integrate it yourself, take a look at [existings linters integrations](https://github.com/golangci/golangci-lint/tree/master/pkg/golinters). Or you can create [GitHub Issue](https://github.com/golangci/golangci-lint/issues/new) and we will integrate it soon.

**Q: It's cool to use `golangci-lint` when starting project, but what about existing projects with large codebase? It will take days to fix all found issues**

A: We are sure that every project can easily integrate `golangci-lint`, event the large one. The idea is to not fix all existing issues. Fix only newly added issue: issues in new code. To do this setup CI (or better use [GolangCI](https://golangci.com) to run `golangci-lint` with option `--new-from-rev=origin/master`. Also take a look at option `-n`.
By doing this you won't create new issues in code and can smoothly fix existing issues (or not).

**Q: How to use `golangci-lint` in CI (Continuous Integration)?**

A: You have 2 choices:
1. Use [GolangCI](https://golangci.com): this service is highly integrated with GitHub (found issues are commented in pull request) and uses `golangci-lint` tool. For configuration use `.golangci.yml` (or toml/json).
2. Use custom CI: just run `golangci-lint` in CI and check exit code. If it's non-zero - fail build. The main disadvantage is that you can't see found issues in pull request code and should view build log, then open needed source file to see a context.

**Q: `golangci-lint` doesn't work**
1. Update it: `go get -u gopkg.in/golangci/golangci-lint.v1/cmd/golangci-lint`
2. Run it with `-v` option and check output.
3. If it doesn't help create [GitHub issue](https://github.com/golangci/golangci-lint/issues/new).
