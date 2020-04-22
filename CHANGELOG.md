Follow the news and releases on our [twitter](https://twitter.com/golangci) and our [blog](https://medium.com/golangci).
There is the most valuable changes log:

## April 2020

1. Add new linters: `testpackage`, `nestif`, `godot`, `gomodguard`
2. Add  github actions output format
3. Update linters: `wsl`
4. Support `disabled-tags` setting for `gocritic`
5. Mitigate OOM and "failed prerequisites"
6. Self-isolate due to unexpected pandemics

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