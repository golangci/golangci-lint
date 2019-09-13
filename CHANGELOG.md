Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).
There is the most valuable changes log:

### June 2019

1. treat Go source files as a plain text by `misspell`: it allows detecting issues in strings, variable names, etc.
2. implement richer and more stable auto-fix of `misspell` issues.

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

### February 2019

1. Implement auto-fixing for `gofmt`, `goimports` and `misspell`
2. Update `unparam`, `goimports`, `gosec` and `go-critic`
3. Support `issues.exclude-rules` config option
4. Add more `identifier` marking patterns
5. Add code-climate output format
6. Fix diff parsing on windows
7. Add version information to built artifact for go1.12
8. Dockerfile: copy the binary to `/usr/bin/` instead of `$GOPATH/bin/`
9. Support `ignore-words` config option for `misspell`
10. Include `staticcheck` check name into a message
11. Fix working with symbolic links

### January 2019

1. Update `megacheck` (`staticcheck`), `unparam` and `go-critic` to the latest versions.
2. Support the new `stylecheck` linter.
3. Support of `enabled-tags` options for `go-critic`.
4. Make rich debugging for `go-critic` and meticulously validate `go-critic` checks config.
5. Update and use upstream versions of `unparam` and `interfacer` instead of forked ones.
6. Improve handling of unknown linter names in `//nolint` directives.
7. Speedup `typecheck` on large project with compilation errors.
8. Add support for searching for `errcheck` exclude file.
9. Fix `go-misc` checksum.
10. Don't crash when staticcheck panics

### December 2018

1. Update `goimports`: the new version creates named imports for name/path mismatches.
2. Update `go-critic` to the latest version.
3. Sync default `go-critic` checks list with the `go-critic`.
4. Support `pre-commit.com` hooks.
5. Rework and simplify `--skip-dirs` for some edge cases.
6. Add `modules-download-mode` option: it's useful in CI.
7. Better validate commands.
8. Fix working with absolute paths.
9. Fix `errcheck.ignore` option.

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