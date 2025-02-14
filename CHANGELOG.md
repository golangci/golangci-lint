Follow the news and releases on [Mastodon](https://fosstodon.org/@golangcilint) and on [Bluesky](https://bsky.app/profile/golangci-lint.run).

`golangci-lint` is a free and open-source project built by volunteers.

If you value it, consider supporting us, we appreciate it! ❤️

[![Open Collective backers and sponsors](https://img.shields.io/badge/OpenCollective-Donate-blue?logo=opencollective&style=for-the-badge)](https://opencollective.com/golangci-lint)
[![GitHub Sponsors](https://img.shields.io/badge/GitHub-Donate-blue?logo=github&style=for-the-badge)](https://github.com/sponsors/golangci)

### v1.64.5

1. Bug fixes
   * Add missing flag `new-from-merge-base-flag`
2. Linters bug fixes
   * `asciicheck`: from 0.3.0 to 0.4.0
   * `forcetypeassert`: from 0.1.0 to 0.2.0
   * `gosec`: from 2.22.0 to 2.22.1

### v1.64.4

1. Linters bug fixes
   * `gci`: fix standard packages list for go1.24

### v1.64.3

1. Linters bug fixes
   * `ginkgolinter`: from 0.18.4 to 0.19.0
   * `go-critic`: from 0.11.5 to 0.12.0
   * `revive`: from 1.6.0 to 1.6.1
   * `gci`: fix standard packages list for go1.24
2. Misc.
   * Build Docker images with go1.24

### v1.64.2

This is the last minor release of golangci-lint v1.
The next release will be golangci-lint [v2](https://github.com/golangci/golangci-lint/issues/5300).

1. Enhancements
   * 🎉 go1.24 support
   * New `issues.new-from-merge-base` option
   * New `run.relative-path-mode` option
2. Linters new features
   * `copyloopvar`: from 1.1.0 to 1.2.1 (support suggested fixes)
   * `exptostd`: from 0.3.1 to 0.4.1 (handles `golang.org/x/exp/constraints.Ordered`)
   * `fatcontext`: from 0.5.3 to 0.7.1 (new option: `check-struct-pointers`)
   * `perfsprint`: from 0.7.1 to 0.8.1 (new options: `integer-format`, `error-format`, `string-format`, `bool-format`, and `hex-format`)
   * `revive`: from 1.5.1 to 1.6.0 (new rules: `redundant-build-tag`, `use-errors-new`. New option `early-return.early-return`)
3. Linters bug fixes
   * `go-errorlint`: from 1.7.0 to 1.7.1
   * `gochecknoglobals`: from 0.2.1 to 0.2.2
   * `godox`: from 006bad1f9d26 to 1.1.0
   * `gosec`: from 2.21.4 to 2.22.0
   * `iface`: from 1.3.0 to 1.3.1
   * `nilnesserr`: from 0.1.1 to 0.1.2
   * `protogetter`: from 0.3.8 to 0.3.9
   * `sloglint`: from 0.7.2 to 0.9.0
   * `spancheck`: fix default `StartSpanMatchersSlice` values
   * `staticcheck`: from 0.5.1 to 0.6.0
4. Deprecations
   * ⚠️ `tenv` is deprecated and replaced by `usetesting.os-setenv: true`.
5. Misc.
   * Sanitize severities by output format
   * Avoid panic with plugin without description
6. Documentation
   * Clarify `depguard` configuration

### v1.64.1

Cancelled due to CI failure.

### v1.64.0

Cancelled due to CI failure.

### v1.63.4

1. Linters bug fixes
   * `dupl`, `gomodguard`, `revive`: keep only Go-files.

### v1.63.3

1. Linters bug fixes
   * `gofmt`, `gofumpt`, `goimports`, `gci`: panic with several trailing EOL
   * `goheader`: skip issues with invalid positions

### v1.63.2

1. Linters bug fixes
   * `gofmt`, `gofumpt`, `goimports`, `gci`: panic with missing trailing EOL

### v1.63.1

1. Linters bug fixes
   * `cgi`: invalid reports with cgo
   * `gofumpt`: panic with autofix and cgo

### v1.63.0

1. Enhancements
   * Add support for SuggestedFixes 🎉 (35 linters can "autofix" reports).
   * Formatters (`gofmt`, `goimports`, `gofumpt`, `gci`) are applied after the suggested fixes.
2. New linters
   * Add `exptostd` linter https://github.com/ldez/exptostd
   * Add `nilnesserr` linter https://github.com/alingse/nilnesserr
   * Add `usetesting` linter https://github.com/ldez/usetesting
3. Linters new features
   * `gci`: new options: `no-inline-comments`, `no-prefix-comments`
   * `gomoddirectives`: from 0.2.4 to 0.6.0 (new options: `go-version-pattern`, `toolchain-pattern`,`toolchain-forbidden`, `tool-forbidden`, `go-debug-forbidden`)
   * `govet`: new `stdversion`, `waitgroup` analyzers
   * `importas`: allow multiple empty aliases
   * `loggercheck`: new `slog` option
   * `recvcheck`: from 0.1.2 to 0.2.0 (new options: `disable-builtin`, `exclusions`)
   * `tagliatelle`: from 0.5.0 to 0.7.1 (new options: `ignored-fields`, `extended-rules`,`overrides`, `pkg`, `ignore`)
   * `usestdlibvars`: from 1.27.0 to 1.28.0 (autofix)
   * `wrapcheck`: from 2.9.0 to 2.10.0 (new option: `extra-ignore-sigs`)
4. Linters bug fixes
   * `asciicheck`: from 0.2.0 to 0.3.0
   * `bodyclose`: from 574207250966 to ed6a65f985e
   * `funlen`: from 0.1.0 to 0.2.0
   * `ginkgolinter`: from 0.18.3 to 0.18.4
   * `gochecksumtype`: from 0.2.0 to 0.3.1
   * `gocognit`: from 1.1.3 to 1.2.0
   * `godot`: from 1.4.18 to 1.4.20
   * `goheader`: report position improvement
   * `gosec`: handling of global nosec option when it is false
   * `iface`: from 1.2.1 to 1.3.0
   * `importas`: from 0.1.0 to 0.2.0
   * `intrange`: from 0.2.1 to 0.3.0
   * `makezero`: from 1.1.1 to 1.2.0
   * `mirror`: from 1.2.0 to 1.3.0
   * `nilnil`: from 1.0.0 to 1.0.1
   * `nosprintfhostport`: from 0.1.1 to 0.2.0
   * `reassign`: from 0.2.0 to 0.3.0
   * `spancheck`: from 0.6.2 to 0.6.4
   * `tagalign`: from 1.3.4 to 1.4.1
   * `wastedassign`: from 2.0.7 to 2.1.0
   * `whitespace`: from 0.1.1 to 0.2.0
   * `wsl`: from 4.4.1 to 4.5.0
5. Deprecations
   * ⚠️ `output.uniq-by-line` is deprecated and replaced by `issues.uniq-by-line`.
6. Misc.
   * Improvements of the help command (color and JSON support).
   * Removes `decoder`, `sloglint`, `tagalign` from `format` preset.
   * Enables paths with junction inside Windows.
   * The timeout is disabled if `run.timeout` <= 0.

### v1.62.2

1. Linters bug fixes
   * `fatcontext`: from 0.5.2 to 0.5.3
   * `ginkgolinter`: from 0.18.0 to 0.18.3
   * `errorlint`: from 1.6.0 to 1.7.0
   * `iface`: from 1.2.0 to 1.2.1
   * `revive`: from 1.5.0 to 1.5.1
   * `testifylint`: from 1.5.0 to 1.5.2
2. Misc.
   * fix: ignore cache error when file not found

### v1.62.1

Cancelled due to CI failure.

### v1.62.0

1. New linters
   * Add `recvcheck` linter https://github.com/raeperd/recvcheck
   * Add `iface` linter https://github.com/uudashr/iface
2. Linters new features
   * `ginkgolinter`: from 0.17.0 to 0.18.0 (new option: `force-succeed`)
   * `gochecksumtype`: from 0.1.4 to 0.2.0 (new option: `default-signifies-exhaustive`)
   * `loggercheck`: from 0.9.4 to 0.10.1 (`log/slog` support)
   * `nilnil`: from 0.1.9 to 1.0.0 (new option: `detect-opposite`)
   * `revive`: from 1.3.9 to 1.5.0 (new rules: `filename-format`, and `file-length-limit`)
   * `tenv`: from 1.10.0 to 1.12.1 (handle dot import)
   * `testifylint`: from 1.4.3 to 1.5.0 (new checkers: `contains`, `encoded-compare`, `regexp`)
3. Linters bug fixes
   * `bidichk`: from 0.2.7 to 0.3.2 (important performance improvement)
   * `canonicalheader`: from 1.1.1 to 1.1.2
   * `cyclop`: from 1.2.1 to 1.2.3
   * `dupword`: from 0.1.1 to 0.1.3
   * `errcheck`: from 1.7.0 to 1.8.0
   * `errchkjson`: from 0.3.6 to 0.4.0
   * `errname`: from 0.1.13 to 1.0.0
   * `gocritic`: from 0.11.4 to 0.11.5
   * `goprintffuncname`: from 7558a9eaa5af to v0.1.0
   * `godot`: from 1.4.17 to 1.4.18
   * `gosec`: from 2.21.2 to 2.21.4
   * `intrange`: from 0.2.0 to 0.2.1
   * `musttag`: from 0.12.2 to 0.13.0
   * `nakedret`: from 2.0.4 to 2.0.5
   * `noctx`: from 0.0.2 to 0.1.0
   * `protogetter`: from 0.3.6 to 0.3.8
4. Deprecations
   * ⚠️ `execinquery`: deprecation step 2
   * ⚠️ `gomnd`: deprecation step 2 (replaced by `mnd`)
5. Misc.
   * Type sizing when cross-compiling (32-bit).
   * code-climate: add check_name field
   * Improve Go version detection
   * Fix Go version propagation
6. Documentation
   * Adds a section about `exclude-dirs-use-default`
   * Improve 'install from sources' section
   * Improve FAQ about Go versions
   * Improve linter/rule/check docs
   * Improve new linter section
   * Improve `forbidigo` pattern examples for built-in functions

### v1.61.0

1. Enhancements
   * Add `junit-xml-extended` format
   * Exclude Swagger Codegen files by default
2. Linters new features
   * `gci`: from 0.13.4 to 0.13.5 (new option `no-lex-order`)
   * `nolintlint`: remove the empty line in the directive replacement
3. Linters bug fixes
   * `dupword`: from 0.0.14 to 0.1.1
   * `fatcontext`: from 0.4.0 to 0.5.2
   * `go-ruleguard`: from 0.4.2 to 0fe6f58b47b1 (fix panic with custom linters)
   * `godot`: from 1.4.16 to 1.4.17
   * `gomodguard`: from 1.3.3 to 1.3.5
   * `gosec`: disable temporarily `G407`
   * `gosec`: from ab3f6c1c83a0 to 2.21.2 (partially fix `G115`)
   * `intrange`: from 0.1.2 to 0.2.0
4. Misc.
    * Improve runtime version parsing
5. Documentation
    * Add additional info about `typecheck`

### v1.60.3

1. Updated linters
   * `gosec`: from 81cda2f91fbe to ab3f6c1c83a0 (fix `G115` false positives)
2. Misc.
   * Check that the Go version use to build is greater or equals to the Go version of the project

### v1.60.2

1. Linters new features
   * `gosec`: from 5f0084eb01a9 to 81cda2f91fbe (adds `G115`, `G405`, `G406`, `G506`, `G507`)
2. Linters bug fixes
   * `gofmt`: update to HEAD (go1.22)
   * `gofumpt`: from 0.6.0 to 0.7.0
   * `gosec`: fix G602 analyzer
   * `staticcheck`: from 0.5.0 to 0.5.1
   * `staticcheck`: propagate Go version
   * `wrapcheck`: from 2.8.3 to 2.9.0
3. Deprecations
   * ⚠️ `exportloopref`: deprecation

### v1.60.1

1. Misc.
   * 🎉 go1.23 support
2. Linters new features
   * `exhaustruct`: from 3.2.0 to 3.3.0 (recognize custom error values in return)
   * `govet` (`printf`): report non-constant format, no args
   * `lll`: advertise max line length instead of just reporting failure
   * `revive`: from 1.3.7 to 1.3.9 (new rule: `comments-density`)
   * `staticcheck`: from 0.4.7 to 0.5.0
   * `testifylint`: from 1.3.1 to 1.4.3 (new options: `formatter`, `suite-broken-parallel`, `suite-subtest-run`)
3. Linters bug fixes
   * `errorlint`: from 1.5.2 to 1.6.0
   * `fatcontext`: from 0.2.2 to 0.4.0 (fix false positives for context stored in structs)
   * `gocognit`: from 1.1.2 to 1.1.3
   * `gomodguard`: from 1.3.2 to 1.3.3
   * `sloglint`: from 0.7.1 to 0.7.2
   * `spancheck`: from 0.6.1 to 0.6.2
   * `tenv`: from 1.7.1 to 1.10.0 (remove reports on fuzzing)
   * `tparallel`: from 0.3.1 to 0.3.2
   * `usestdlibvars`: from 1.26.0 to 1.27.0 (fix false-positive with number used inside a mathematical operations)
   * `wsl`: from 4.2.1 to 4.4.1
4. Deprecations
   * ️⚠️ `unused`: remove `exported-is-used` option
5. Fixes
   * SARIF: sanitize level property
   * ️⚠️ `typecheck` issues should never be ignored
6. Documentation
   * Add link on linter without configuration
   * Remove 'trusted by' page
   * `wsl` update documentation of the configuration

### v1.60.0

Cancelled due to a CI problem.

### v1.59.1

1. Updated linters
   * `errorlint`: from 1.5.1 to 1.5.2
   * `gomnd`: deprecated configuration compatibility
   * `intrange`: add `style` preset
   * `misspell`: from 0.5.1 to 0.6.0
   * `sloglint`: from 0.7.0 to 0.7.1
   * `testifylint`: from 1.3.0 to 1.3.1
   * `unparam`: bump to HEAD
   * `usestdlibvars`: from 1.25.0 to 1.26.0
2. Fixes
   * SARIF: init empty result slice
   * SARIF: issue column >= 1
3. Documentation
   * `revive`: update documentation of the configuration

### v1.59.0

1. Enhancements
   * Add SARIF output format
   * Allow the analysis of generated files (`issues.exclude-generated: disable`)
2. Updated linters
   * `errcheck`: fix deprecation warning
   * `gocritic`: from 0.11.3 to 0.11.4
   * `gosec`: from 2.20.0 to 5f0084eb01a9 (fix G601 and G113 performance issues)
   * `sloglint`: from 0.6.0 to 0.7.0 (new option `forbidden-keys`)
   * `testifylint`: from 1.2.0 to 1.3.0 (new checker `negative-positive` and new option `go-require.ignore-http-handlers`)
3. Misc.
   * ️️⚠️ Deprecate `github-action` output format
   * ️️⚠️ Deprecate `issues.exclude-generated-strict` option (replaced by `issues.exclude-generated: strict`)
   * ️️⚠️ Add warning about disabled and deprecated linters (level 2)

### v1.58.2

1. Updated linters
   * `canonicalheader`: from 1.0.6 to 1.1.1
   * `gosec`: from 2.19.0 to 2.20.0
   * `musttag`: from 0.12.1 to 0.12.2
   * `nilnil`: from 0.1.8 to 0.1.9
2. Documentation
   * Improve integrations and install pages

### v1.58.1

1. Updated linters
   * `tagalign`: from 1.3.3 to 1.3.4
   * `protogetter`: from 0.3.5 to 0.3.6
   * `gochecknoinits`: fix analyzer name
2. Fixes
   * Restores previous `github-actions` output format (removes GitHub Action problem matchers)

### v1.58.0

1. New linters
   * `fatcontext`: https://github.com/Crocmagnon/fatcontext
   * `canonicalheader`: https://github.com/lasiar/canonicalheader
2. Updated linters
   * `copyloopvar`: from 1.0.10 to 1.1.0 (`ignore-alias` is replaced by `check-alias` with the opposite behavior)
   * `decorder`: from 0.4.1 to 0.4.2
   * `errname`: from 0.1.12 to 0.1.13
   * `errorlint`: from 1.4.8 to 1.5.1 (new options `allowed-errors` and `allowed-errors-wildcard`)
   * `execinquery`: deprecate linter ⚠️
   * `gci`: from 0.12.3 to 0.13.4 (new section `localModule`)
   * `gocritic`: from 0.11.2 to 0.11.3
   * `spancheck`: from 0.5.3 to 0.6.1
   * `goerr113` is replaced by `err113` ⚠️
   * `gomnd` is replaced by `mnd` ⚠️
   * `gomodguard`: from 1.3.1 to 1.3.2
   * `grouper`: from 1.1.1 to 1.1.2
   * `intrange`: from 0.1.1 to 0.1.2
   * `mirror`: from 1.1.0 to 1.2.0
   * `misspell`: from 0.4.1 to 0.5.1
   * `musttag`: from 0.9.0 to 0.12.1
   * `nilnil`: from 0.1.7 to 0.1.8
   * `nonamedreturns`: from 1.0.4 to 1.0.5
   * `promlinter`: from 0.2.0 to 0.3.0
   * `sloglint`: from 0.5.0 to 0.6.0
   * `unparam`: bump to HEAD (063aff900ca150b80930c8de76f11d7e6488222f)
   * `whitespace`: from 0.1.0 to 0.1.1
3. Enhancements
   * Speed up "fast" linters when only "fast" linters are run: between 40% and 80% faster at first run (i.e. without cache)
4. Fixes
   * Use version with module plugins
   * Skip `go.mod` report inside autogenerated processor
   * Keep only `typecheck` issues when needed
   * Don't hide `typecheck` errors inside diff processor
5. Misc.
   * ⚠️ log an error when using previously deprecated linters ([Linter Deprecation Cycle](https://golangci-lint.run/product/roadmap/#linter-deprecation-cycle))
      * [`deadcode`](https://github.com/remyoudompheng/go-misc/tree/HEAD/deadcode): deprecated since v1.49.0 (2022-08-23).
      * [`exhaustivestruct`](https://github.com/mbilski/exhaustivestruct): deprecated since v1.46.0 (2022-05-08).
      * [`golint`](https://github.com/golang/lint): deprecated since v1.41.0 (2021-06-15).
      * [`ifshort`](https://github.com/esimonov/ifshort): deprecated since v1.48.0 (2022-08-04).
      * [`interfacer`](https://github.com/mvdan/interfacer): deprecated since v1.38.0 (2021-03-03).
      * [`maligned`](https://github.com/mdempsky/maligned): deprecated since v1.38.0 (2021-03-03).
      * [`nosnakecase`](https://github.com/sivchari/nosnakecase): deprecated since v1.48.0 (2022-08-04).
      * [`scopelint`](https://github.com/kyoh86/scopelint): deprecated since v1.39.0 (2021-03-25).
      * [`structcheck`](https://github.com/opennota/check): deprecated since v1.49.0 (2022-08-23).
      * [`varcheck`](https://github.com/opennota/check): deprecated since v1.49.0 (2022-08-23).
   * ⚠️ Deprecate usage of linter alternative names
   * Remove help display on errors with `config verify` command
   * Add `pre-commit` hook to run `config verify`
   * Improve `github-action` output
6. Documentation
   * Remove deprecated Atom from Editor Integrations

GitHub Action (v5.1.0) for golangci-lint:
- supports for `pull`, `pull_request_target`, and `merge_group` events with the option `only-new-issues`.
- ️️⚠️ `skip-pkg-cache` and `skip-build-cache` have been removed because the cache related to Go itself is already handled by `actions/setup-go`.
- with golangci-lint v1.58, the file information (path and position) will be displayed on the log.

### v1.57.2

1. Updated linters
   * `contextcheck`: from 1.1.4 to 1.1.5
   * `copyloopvar`: from 1.0.8 to 1.0.10
   * `ginkgolinter`: from 0.16.1 to 0.16.2
   * `goconst`: from 1.7.0 to 1.7.1
   * `gomoddirectives`: from 0.2.3 to 0.2.4
   * `intrange`: from 0.1.0 to 0.1.1
2. Misc.
   * Display warnings on deprecated linter options
   * Fix missing `colored-tab` output format
   * Fix TeamCity `inspectionType` service message
3. Documentation
   * Remove invalid example about mixing files and directory
   * Improve linters page

### v1.57.1

1. Fixes
   * Ignore issues with invalid position (e.g. `contextcheck`).

### v1.57.0

1. New linters
   * `copyloopvar`: https://github.com/karamaru-alpha/copyloopvar
   * `intrange`: https://github.com/ckaznocha/intrange
2. Updated linters
   * `dupword`: from 0.0.13 to 0.0.14
   * `gci`: from 0.12.1 to 0.12.3
   * `ginkgolinter`: from 0.15.2 to 0.16.1 (new option `force-expect-to`, `validate-async-intervals`, and `forbid-spec-pollution`)
   * `go-critic`: from 0.11.1 to 0.11.2
   * `go-critic`: support of `enable-all` and `disable-all` options
   * `go-spancheck`: from 0.5.2 to 0.5.3
   * `gomodguard`: from 1.3.0 to 1.3.1
   * `govet`: deprecation of `check-shadowing` ⚠️
   * `govet`: disable temporarily `httpresponse` because of a bug https://github.com/golang/go/issues/66259
   * `misspell`: add `extra-words`
   * `musttag`: from 0.8.0 to 0.9.0
   * `nakedret`: from 2.0.2 to 2.0.4
   * `paralleltest`: from 1.0.9 to 1.0.10
   * `perfsprint`: from 0.6.0 to 0.7.1 (new option `strconcat`)
   * `protogetter`: from 0.3.4 to 0.3.5
   * `revive`: add `exclude` option
   * `sloglint`: from 0.4.0 to 0.5.0 (new option `no-global`)
   * `staticcheck`: from 0.4.6 to 0.4.7
   * `testifylint`: from 1.1.2 to 1.2.0 (new option `bool-compare`)
   * `unconvert`: to HEAD (new options `fast-math` and `safe`)
   * `wrapcheck`: from 2.8.1 to 2.8.3
   * Disable `copyloopvar` and `intrange` on Go < 1.22
3. Enhancements
   * 🧩New custom linters system https://golangci-lint.run/plugins/module-plugins/
   * Allow running only a specific linter without modifying the file configuration (`--enable-only`)
   * Allow custom sort order for the reports (`output.sort-order`)
   * Automatically adjust the maximum concurrency to the container CPU quota if `run.concurrency=0`
   * Add `config verify` command to check the configuration against the JSON Schema
   * Option to strictly follow Go generated file convention (`issues.exclude-generated-strict`)
   * Syntax to not override `severity` from linters (`@linter`)
   * Use severities from `gosec`
   * Create automatically directory related to `output.formats.path`
   * Use the first issue without inline on `mergeLineIssues` on multiple issues
4. Misc.
   * ⚠️ Inactivate deprecated linters (`deadcode`, `exhaustivestruct`, `golint`, `ifshort`, `interfacer`, `maligned`, `nosnakecase`, `scopelint`, `structcheck`, `varcheck`)
   * ⚠️ Deprecated CLI flags have been removed (deprecated since 2018)
   * ⚠️ Move `show-stats` option from `run` to `output` configuration section
   * ⚠️ Replace `run.skip-xxx` options by `issues.exclude-xxx` options
   * ⚠️ Replace `output.format` by `output.formats` with a new file configuration syntax
   * Internal rewrite of the CLI
   * Improve 'no go files to analyze' message
   * Use `GOTOOLCHAIN=auto` inside the Docker images
5. Documentation
   * ⚠️ Define the linter deprecation cycle https://golangci-lint.run/product/roadmap/#linter-deprecation-cycle
   * 🎉Use information from the previous release to create linter pages
   * Publish JSON schema on https://golangci-lint.run/jsonschema/golangci.jsonschema.json
   * Reorganize documentation pages
   * Add an explanation about the configuration file inside golangci-lint repository

**⚠️ Important ⚠️**

1. Deprecated linters are inactivated, you still need to disable them if you are using `enable-all`.
2. Deprecated CLI flags (about linter settings and `deadline`) have been removed.

### v1.56.2

1. updated linters
   * `go-critic`: from 0.11.0 to 0.11.1
   * `gosec`: from 2.18.2 to 2.19.0
   * `testifylint`: from 1.1.1 to 1.1.2
   * `usestdlibvars`: from 1.24.0 to 1.25.0
   * `wsl`: from 4.2.0 to 4.2.1
2. misc.
   * Fix missing version in Docker image
3. Documentation
   * Explain the limitation of `new-from-rev` and `new-from-patch`

### v1.56.1

1. updated linters
   * `errcheck`: from 1.6.3 to 1.7.0
   * `govet`: disable `loopclosure` with go1.22
   * `revive`: from 1.3.6 to 1.3.7
   * `testifylint`: from 1.1.0 to 1.1.1

### v1.56.0

1. new linters
   * `spancheck`: https://github.com/jjti/go-spancheck
2. updated linters
   * `depguard`: from 2.1.0 to 2.2.0
   * `exhaustive`: from 0.11.0 to 0.12.0
   * `exhaustruct`: from 3.1.0 to 3.2.0
   * `gci`: from 0.11.2 to 0.12.1
   * `ginkgolinter`: from 0.14.1 to 0.15.2
   * `go-check-sumtype`: from 0.1.3 to 0.1.4
   * `go-critic`: from 0.9.0 to 0.11.0
   * `go-errorlint`: from 1.4.5 to 1.4.8
   * `go-spancheck`: from 0.4.2 to 0.5.2
   * `goconst`: from 1.6.0 to 1.7.0
   * `godot`: from 1.4.15 to 1.4.16
   * `gofumpt`: from 0.5.0 to 0.6.0
   * `inamedparam`: from 0.1.2 to 0.1.3
   * `ineffassign`: from 0.0.0-20230610083614-0e73809eb601 to 0.1.0
   * `ireturn`: from 0.2.2 to 0.3.0
   * `misspell`: add mode option
   * `musttag`: from v0.7.2 to v0.8.0
   * `paralleltest`: from 1.0.8 to 1.0.9
   * `perfsprint`: from 0.2.0 to 0.6.0
   * `protogetter`: from 0.2.3 to 0.3.4
   * `revive`: from 1.3.4 to 1.3.6
   * `sloglint`: add static-msg option
   * `sloglint`: from 0.1.2 to 0.4.0
   * `testifylint`: from 0.2.3 to 1.1.0
   * `unparam`: from 20221223090309-7455f1af531d to 20240104100049-c549a3470d14
   * `whitespace`: update after moving to the `analysis` package
   * `wsl`: from 3.4.0 to 4.2.0
   * `zerologlint`: from 0.1.3 to 0.1.5
3. misc.
   * 🎉 go1.22 support
   * Implement stats per linter with a flag
   * Make versioning inside Docker image consistent with binaries
   * Parse Go RC version
4. Documentation
   * Fix `noctx` description
   * Add missing fields to `.golangci.reference.yml`
   * Improve `.golangci.reference.yml` defaults
   * `typecheck`: improve FAQ
   * `exhaustruct`: note that struct regular expressions are expected to match the entire `package/name/structname`
   * `wrapcheck`: adjust `ignoreSigs` to new defaults

**Important**

`testifylint` has [breaking changes](https://github.com/Antonboom/testifylint/releases/tag/v1.0.0) about enabling/disabling checks:

- If you were using the option `enable` with a filtered list of checks, you should either add `disable-all: true` (1) or use `disable` field (2).

    ```yml
      # Example (1)
      testifylint:
        disable-all: true
        enable:
          - bool-compare
          - compares
          - empty
          - error-is-as
          - error-nil
          - expected-actual
          - go-require
          - float-compare
          - len
          - nil-compare
          - require-error
          # - suite-dont-use-pkg
          - suite-extra-assert-call
          - suite-thelper
    ```

    ```yml
      # Example (2)
      testifylint:
        disable:
          - suite-dont-use-pkg
    ```

### v1.55.2

1. updated linters
   * `ireturn`: from 0.2.1 to 0.2.2
   * `ginkgolinter`: from 0.14.0 to 0.14.1

### v1.55.1

1. updated linters
   * `gosec`: from 2.18.1 to 2.18.2
2. misc.
   * `revgrep`: from v0.5.0 to v0.5.2 (support git < 2.41.0)
   * output: convert backslashes to forward slashes for GitHub Action annotations printer

### v1.55.0

1. new linters
   * `gochecksumtype`: https://github.com/alecthomas/go-check-sumtype
   * `inamedparam`: https://github.com/macabu/inamedparam
   * `perfsprint`: https://github.com/catenacyber/perfsprint
   * `protogetter`: https://github.com/ghostiam/protogetter
   * `sloglint`: https://github.com/go-simpler/sloglint
   * `testifylint`: https://github.com/Antonboom/testifylint
2. updated linters
   * `bidichk`: from 0.2.4 to 0.2.7
   * `decorder`: from 0.4.0 to 0.4.1
   * `dupword`: from 0.0.12 to 0.0.13
   * `errchkjson`: from 0.3.1 to 0.3.6
   * `gci`: from 0.11.0 to 0.11.2
   * `ginkgolinter`: from 0.13.5 to 0.14.0
   * `go-errorlint`: from 1.4.4 to 1.4.5
   * `gocognit`: from 1.0.7 to 1.1.0
   * `goconst`: from 1.5.1 to 1.6.0
   * `godot`: from 1.4.14 to 1.4.15
   * `gofmt`: update to HEAD
   * `goimports`: update to HEAD
   * `gosec`: from 2.17.0 to 2.18.1
   * `gosmopolitan`: from 1.2.1 to 1.2.2
   * `govet`: add `appends` analyzer
   * `ireturn`: from 0.2.0 to 0.2.1
   * `protogetter`: from 0.2.2 to 0.2.3
   * `revgrep`: from 745bb2f7c2e6 to v0.5.0
   * `revive`: from 1.3.2 to 1.3.4
   * `sqlclosecheck`: from 0.4.0 to 0.5.1
   * `staticcheck`: from 0.4.5 to 0.4.6
   * `tagalign`: from 1.3.2 to 1.3.3
   * `unused`: support passing in options
3. misc.
   * Add a pre-commit hook to check all files
4. Documentation
   * add source options to exclude-rules docs
   * `gosec`: add G602 to includes/excludes inside .golangci.reference.yml

### v1.54.2

1. updated linters:
   * `errname`: from 0.1.10 to 0.1.12
   * `ginkgolinter`: from 0.13.3 to 0.13.5
   * `go-errorlint`: from 1.4.3 to 1.4.4
   * `godot`: from 1.4.11 to 1.4.14
   * `gosec`: from 2.16.0 to 2.17.0
   * `musttag`: from 0.7.1 to 0.7.2
   * `nilnil`: from 0.1.5 to 0.1.7
   * `staticcheck`: from 0.4.3 to 0.4.5
   * `usestdlibvars`: from 1.23.0 to 1.24.0
   * `govet`: add missing `directive` and `slog` passes

### v1.54.1

1. updated linters:
   * `go-critic`:  from 0.8.2 to 0.9.0
2. misc.
   * plugin: temporarily hide warning about using plugins using the old API

### v1.54.0

1. updated linters:
   * `decorder`: from 0.2.3 to 0.4.0
   * `dupword`: from 0.0.11 to 0.0.12
   * `errorlint`: from 1.4.2 to 1.4.3
   * `exhaustruct`: from 2.3.0 to 3.1.0
   * `forbidigo`: from 1.5.3 to 1.6.0
   * `funlen`: from 0.0.3 to 0.1.0
   * `gci`: from 0.10.1 to 0.11.0
   * `ginkgolinter`: from 0.12.1 to 0.13.3
   * `go-critic`: from 0.8.1 to 0.8.2
   * `go-errorlint`: from 1.4.2 to 1.4.3
   * `go-exhaustruct`: from 2.3.0 to 3.1.0
   * `gocognit`: from 1.0.6 to 1.0.7
   * `gocritic`: from 0.8.1 to 0.8.2
   * `gofmt`: autofix missing newline at EOF
   * `misspell`: 0.4.0 to 0.4.1
   * `musttag`: from 0.7.0 to 0.7.1
   * `paralleltest`: from 1.0.7 to 1.0.8
   * `tagalign`: from 1.2.2 to 1.3.2
   * `typecheck`: explain it and remove it from the linter list
   * `zerologlint`: from 0.1.2 to 0.1.3
2. misc.
   * 🎉 go1.21 support
   * plugin: include custom linters in `enable-all`
   * plugin: allow to use settings for plugins
3. Documentation
   * Add linter descriptions.

**Important**

`ruleguard` (a "rule" inside `gocritic`) was disabled in this release (v1.54.0) and was enabled again in the next release (v1.54.1).

`exhaustruct` has breaking changes with regular expressions, more details [here](https://github.com/GaijinEntertainment/go-exhaustruct/releases/tag/v3.0.0).

### v1.53.3

1. updated linters:
   * `cyclop`: from 1.2.0 to 1.2.1
   * `exhaustive`: from 0.10.0 to 0.11.0
   * `forbidigo`: from 1.5.1 to 1.5.3
   * `ginkgolinter`: from 0.12.2 to 0.12.1
   * `ineffassign`: bump to HEAD
   * `nakedret`: from 2.0.1 to 2.0.2
   * `zerologlint`: from 0.1.1 to 0.1.2
2. misc.
   * codeclimate: reduce memory allocation
   * support illumos/amd64

### v1.53.2

1. updated linters
   * `depguard`: from v2.0.1 to 2.1.0
2. misc.
   * `depguard`: throw error only when the linter is called

### v1.53.1

1. misc.
   * `depguard`: fix GOROOT detection
   * `importas`: fix duplication detection when aliases use regular expression replacement pattern

### v1.53.0

1. new linters
   * `gosmopolitan`: https://github.com/xen0n/gosmopolitan
   * `mirror`: https://github.com/butuzov/mirror
   * `tagalign`: https://github.com/4meepo/tagalign
   * `zerologlint`: https://github.com/ykadowak/zerologlint
2. updated linters
   * `bodyclose`: bump to HEAD
   * `containedctx`: from 1.0.2 to 1.0.3
   * `depguard`: migrate to v2
   * `errname`: from 0.1.9 to 0.1.10
   * `exhaustive`: from 0.9.5 to 0.10.0
   * `forbidigo`: better support for configuring complex rules
   * `gci`: improve error message
   * `ginkgolinter`: add suppress-async-assertion option
   * `ginkgolinter`: from 0.9.0 to 0.12.0
   * `go-critic`: from 0.7.0 to 0.8.1
   * `go-errorlint`: from 1.4.0 to 1.4.2
   * `gofumpt`: from 0.4.0 to 0.5.0
   * `gosec`: convert global settings as map with proper key type
   * `gosec`: from 2.15.0 to 2.16.0
   * `importas`: detect duplicate alias or package in the configuration
   * `importas`: fix typo in logs
   * `ireturn`: from 0.1.1 to 0.2.0
   * `musttag`: from 0.5.0 to 0.7.0
   * `nakedret`: to 2.0.1
   * `nilnil`: from 0.1.3 to 0.1.5
   * `paralleltest`: from 1.0.6 to 1.0.7
   * `revive`: from 1.3.1 to 1.3.2
   * `tagliatelle`: from 0.4.0 to 0.5.0
   * `usestdlibvars`: fix configuration
3. misc.
   * `golang.org/x/tools`: from 0.7.0 to 0.9.2
   * add loongarch64 to the install script
   * output: add colored-tab
   * remove warning when the config file is explicitly stdin
   * rules: support inverted path match
4. Documentation
   * `mnd`: clarify ignore usage examples to be regexps
   * `tagliatelle`: update documentation
   * improve features section
   * update supported Go versions FAQ

### v1.52.2

1. updated linters
   * `tparallel`: from 0.3.0 to 0.3.1
2. misc.
   * fix: pre-commit `require_serial` and `pass_filenames`

### v1.52.1

1. misc.
   * fix: improve panic management
   * fix: the combination of --fix and --path-prefix

### v1.52.0

1. updated linters
   * `asciicheck`: from 0.1.1 to 0.2.0
   * `bidichk`: from 0.2.3 to 0.2.4
   * `contextcheck`: from 1.1.3 to 1.1.4
   * `dupword`: from 0.0.9 to 0.0.11
   * `durationcheck`: from 0.0.9 to 0.0.10
   * `errchkjson`: from 0.3.0 to 0.3.1
   * `errname`: from 0.1.7 to 0.1.9
   * `forbidigo`: from 1.4.0 to 1.5.1
   * `gci`: from 0.9.1 to 0.10.1
   * `ginkgolinter`: from 0.8.1 to 0.9.0
   * `go-critic`: from 0.6.7 to 0.7.0
   * `go-errorlint`: from 1.1.0  to 1.4.0
   * `godox`: bump to HEAD
   * `lll`: skip go command
   * `loggercheck`: from 0.9.3 to 0.9.4
   * `musttag`: from 0.4.5 to 0.5.0
   * `nilnil`: from 0.1.1 to 0.1.3
   * `noctx`: from 0.0.1 to 0.0.2
   * `revive`: from 1.2.5 to 1.3.1
   * `rowserrcheck`: remove limitation related to generics support
   * `staticcheck`: from 0.4.2 to 0.4.3
   * `testpackage`: from 1.1.0 to 1.1.1
   * `tparallel`: from 0.2.1 to 0.3.0
   * `wastedassign`: remove limitation related to generics support
   * `wrapcheck`: from 2.8.0 to 2.8.1
2. misc.
   * Add TeamCity output format
   * Consider path prefix when matching path patterns
   * Add Go version to version information
3. Documentation
   * Add Tekton in Trusted By page
   * Clarify that custom linters are not enabled by default
   * Remove description for deprecated "go" option

### v1.51.2

1. updated linters
   * `forbidigo`: from 1.3.0 to 1.4.0
   * `gci`: from 0.9.0 to 0.9.1
   * `go-critic`: from 0.6.5 to 0.6.7
   * `go-errorlint`: from 1.0.6 to 1.1.0
   * `gosec`: from 2.14.0 to 2.15.0
   * `musttag`: from 0.4.4 to 0.4.5
   * `staticcheck`: from 0.4.0 to 0.4.2
   * `tools`: from 0.5.0 to 0.6.0
   * `usestdlibvars`: from 1.21.1 to 1.23.0
   * `wsl`: from 3.3.0 to 3.4.0
   * `govet`: enable `timeformat` by default
2. misc.
   * fix: cache status size calculation
   * add new source archive
3. Documentation
   * Improve installation section
   * Replace links to godoc.org with pkg.go.dev

### v1.51.1

1. updated linters
   * `ginkgolinter`: from 0.7.1 to 0.8.1
   * `ineffassign`: bump to HEAD
   * `musttag`: from 0.4.3 to 0.4.4
   * `sqlclosecheck`: from 0.3.0 to 0.4.0
   * `staticcheck`: bump to v0.4.0
   * `wastedassign`: from 2.0.6 to 2.0.7
   * `wrapcheck`: from 2.7.0 to 2.8.0

### v1.51.0

1. new linters
   * `ginkgolinter`: https://github.com/nunnatsa/ginkgolinter
   * `musttag`: https://github.com/tmzane/musttag
   * `gocheckcompilerdirectives`: https://github.com/leighmcculloch/gocheckcompilerdirectives
2. updated linters
   * `bodyclose`: to HEAD
   * `dupword`: from 0.0.7 to 0.0.9
   * `errcheck`: from 1.6.2 to 1.6.3
   * `exhaustive`: from 0.8.3 to 0.9.5
   * `exportloopref`: from 0.1.8 to 0.1.11
   * `gci`: from 0.8.1 to 0.9.0
   * `ginkgolinter`: from 0.6.0 to 0.7.1
   * `go-errorlint`: from 1.0.5 to 1.0.6
   * `go-ruleguard`: from 0.3.21 to 0.3.22
   * `gocheckcompilerdirectives`: from 1.1.0 to 1.2.1
   * `gochecknoglobals`: from 0.1.0 to 0.2.1
   * `gomodguard`: from 1.2.4 to 1.3.0
   * `gosec`: from 2.13.1 to 2.14.0
   * `govet`: Add `timeformat` to analysers
   * `grouper`: from 1.1.0 to 1.1.1
   * `musttag`: from 0.4.1 to 0.4.3
   * `revive`: from 1.2.4 to 1.2.5
   * `tagliatelle`: from 0.3.1 to 0.4.0
   * `tenv`: from 1.7.0 to 1.7.1
   * `unparam`: bump to HEAD
   * `usestdlibvars`: from 1.20.0 to 1.21.1
   * `wsl`: fix `force-err-cuddling` flag
3. misc.
   * go1.20 support
   * remove deprecated linters from presets
   * Build NetBSD binaries
   * Build loong64 binaries
4. Documentation
   * `goimport`: improve documentation for local-prefixes
   * `gomnd`: add missing always ignored functions
   * `nolint`: fix typo
   * `tagliatelle` usage typo
   * add note about binary requirement for plugin
   * cache preserving and colored output on docker runs
   * improve documentation about debugging.
   * improve Editor Integration section
   * More specific default cache directory
   * update output example to use valid checkstyle example; add json example

### v1.50.1

1. updated linters
   * `contextcheck`: from 1.1.2 to 1.1.3
   * `go-mnd`: from 2.5.0 to 2.5.1
   * `wrapcheck`: from 2.6.2 to 2.7.0
   * `revive`: fix configuration parsing
   * `lll`: skip imports
2. misc.
   * windows: remove redundant character escape '\/'
   * code-climate: add default severity

### v1.50.0

1. new linters
   * `dupword`: https://github.com/Abirdcfly/dupword
   * `testableexamples`: https://github.com/maratori/testableexamples
2. updated linters
   * `contextcheck`: change owner
   * `contextcheck`: from 1.0.6 to 1.1.2
   * `depguard`: from 1.1.0 to 1.1.1
   * `exhaustive`: add missing config
   * `exhaustive`: from 0.8.1 to 0.8.3
   * `gci`: from 0.6.3 to 0.8.0
   * `go-critic`: from 0.6.4 to 0.6.5
   * `go-errorlint`: from 1.0.2 to 1.0.5
   * `go-reassign`: v0.1.2 to v0.2.0
   * `gofmt`: add option `rewrite-rules`
   * `gofumpt` from 0.3.1 to 0.4.0
   * `goimports`: update to HEAD
   * `interfacebloat`: fix configuration loading
   * `logrlint`: rename `logrlint` to `loggercheck`
   * `paralleltest`: add tests of the ignore-missing option
   * `revive`: from 1.2.3 to 1.2.4
   * `usestdlibvars`: from 1.13.0 to 1.20.0
   * `wsl`: support all configs and update docs
3. misc.
   * Normalize `exclude-rules` paths for Windows
   * add riscv64 to the install script
4. Documentation
   * cli: remove reference to old service

### v1.49.0

IMPORTANT: `varcheck` and `deadcode` has been removed of default linters.

1. new linters
   * `interfacebloat`: https://github.com/sashamelentyev/interfacebloat
   * `logrlint`: https://github.com/timonwong/logrlint
   * `reassign`: https://github.com/curioswitch/go-reassign
2. updated linters
   * `go-colorable`: from 0.1.12 to 0.1.13
   * `go-critic`: from 0.6.3 to 0.6.4
   * `go-errorlint`: from 1.0.0 to 1.0.2
   * `go-exhaustruct`: from 2.2.2 to 2.3.0
   * `gopsutil`: from 3.22.6 to 3.22.7
   * `gosec`: from 2.12.0 to 2.13.1
   * `revive`: from 1.2.1 to 1.2.3
   * `usestdlibvars`: from 1.8.0 to 1.13.0
   * `contextcheck`: from v1.0.4 to v1.0.6 && re-enable
   * `nosnakecase`: This linter is deprecated.
   * `varcheck`: This linter is deprecated use `unused` instead.
   * `deadcode`: This linter is deprecated use `unused` instead.
   * `structcheck`: This linter is deprecated use `unused` instead.
3. documentation
   * `revive`: fix wrong URL
   * Add a section about default exclusions
   * `usestdlibvars`: fix typo in documentation
   * `nolintlint`: remove allow-leading-space option
   * Update documentation and assets
4. misc.
   * dev: rewrite the internal tests framework
   * fix: exit early on run --version
   * fix: set an explicit `GOROOT` in the Docker image for `go-critic`

### v1.48.0

1. new linters
   * `usestdlibvars`:https://github.com/sashamelentyev/usestdlibvars
2. updated linters
   * `contextcheck`: disable linter
   * `errcheck`: from 1.6.1 to 1.6.2
   * `gci`: add missing `custom-order` setting
   * `gci`: from 0.5.0 to 0.6.0
   * `ifshort`: deprecate linter
   * `nolint`: drop allow-leading-space option and add "nolint:all"
   * `revgrep`: bump to HEAD
3. documentation
   * remove outdated info on source install
4. misc
   * go1.19 support

### v1.47.3

1. updated linters:
   * remove some go1.18 limitations
   * `asasalint`: from 0.0.10 to 0.0.11
   * `decorder`: from 0.2.2 to v0.2.3
   * `gci`: fix panic with invalid configuration option
   * `gci`: from 0.4.3 to v0.5.0
   * `go-exhaustruct`: from 2.2.0 to 2.2.2
   * `gomodguard`: from 1.2.3 to 1.2.4
   * `nosnakecase`: from 1.5.0 to 1.7.0
   * `honnef.co/go/tools`: from 0.3.2 to v0.3.3
2. misc
   * cgo: fix linters ignoring CGo files

### v1.47.2

1. updated linters:
   * `revive`: ignore slow rules

### v1.47.1

1. updated linters:
   * `gci`: from 0.4.2 to 0.4.3
   * `gci`: remove the use of stdin
   * `gci`: fix options display
   * `tenv`: from 1.6.0 to 1.7.0
   * `unparam`: bump to HEAD

### v1.47.0

1. new linters:
   * `asasalint`: https://github.com/alingse/asasalint
   * `nosnakecase`: https://github.com/sivchari/nosnakecase
2. updated linters:
   * `decorder`: from 0.2.1 to 0.2.2
   * `errcheck`: from 1.6.0 to 1.6.1
   * `errname`: from 0.1.6 to 0.1.7
   * `exhaustive`: from 0.7.11 to 0.8.1
   * `gci`: fix issues and re-enable autofix
   * `gci`: from 0.3.4 to 0.4.2
   * `go-exhaustruct`: from 2.1.0 to 2.2.0
   * `go-ruleguard`: from 0.3.19 to 0.3.21
   * `gocognit`: from 1.0.5 to 1.0.6
   * `gocyclo`: from 0.5.1 to 0.6.0
   * `golang.org/x/tools`: bump to HEAD
   * `gosec`: allow `global` config
   * `gosec`: from 2.11.0 to 2.12.0
   * `nonamedreturns`: from 1.0.1 to 1.0.4
   * `paralleltest`: from 1.0.3 to 1.0.6
   * `staticcheck`: fix generics
   * `staticcheck`: from 0.3.1 to 0.3.2
   * `tenv`: from 1.5.0 to 1.6.0
   * `testpackage`: from 1.0.1 to 1.1.0
   * `thelper`: from 0.6.2 to 0.6.3
   * `wrapcheck`: from 2.6.1 to 2.6.2
3. documentation:
   * add thanks page
   * add a clear explanation about the `staticcheck` integration.
   * `depguard`: add `ignore-file-rules`
   * `depguard`: adjust phrasing
   * `gocritic`: add `enable` and `disable` ruleguard settings
   * `gomnd`: fix typo
   * `gosec`: add configs for all existing rules
   * `govet`: add settings for `shadow` and `unusedresult`
   * `thelper`: add `fuzz` config and description
   * linters: add defaults

### v1.46.2

1. updated linters:
   * `execinquery`: bump from v1.2.0 to v1.2.1
   * `errorlint`: bump to v1.0.0
   * `thelper`: allow to disable one option
2. documentation:
   * rename `.golangci.example.yml` to `.golangci.reference.yml`
   * add `containedctx` linter to the list of available linters

### v1.46.1

1. updated linters:
   * `execinquery`: bump from v0.6.0 to v0.6.1
2. documentation:
   * add missing linters

### v1.46.0

1. new linters:
   * `execinquery`: https://github.com/lufeee/execinquery
   * `nonamedreturns`: https://github.com/firefart/nonamedreturns
   * `nosprintfhostport`: https://github.com/stbenjam/no-sprintf-host-port
   * `exhaustruct`: https://github.com/GaijinEntertainment/go-exhaustruct
2. updated linters:
   * `bidichk`: from 0.2.2 to 0.2.3
   * `deadcode`: bump to HEAD
   * `errchkjson`: from 0.2.3 to 0.3.0
   * `errname`: from 0.1.5 to 0.1.6
   * `go-critic`: from 0.6.2 to 0.6.3
   * `gocyclo`: from 0.4.0 to 0.5.1
   * `gofumpt` from 0.3.0 to 0.3.1
   * `gomoddirectives`: from 0.2.2 to 0.2.3
   * `gosec`: from 2.10.0 to 2.11.0
   * `honnef.co/go/tools`: from 0.2.2to 0.3.1 (go1.18 support)
   * `nilnil`: from 0.1.0 to 0.1.1
   * `nonamedreturns`: bump from 1.0.0 to 1.0.1
   * `predeclared`: from 0.2.1 to 0.2.2
   * `promlinter`: bump to v0.2.0
   * `revive`: from 1.1.4 to 1.2.1
   * `tenv`: from 1.4.7 to 1.5.0
   * `thelper`: from 0.5.1 to 0.6.2
   * `unused`: fix false-positive
   * `varnamelen`: bump to v0.8.0
   * `wrapcheck`: from 2.5.0 to 2.6.1
   * `exhaustivestruct`: This linter is deprecated use `exhaustruct` instead.
3. documentation:
   * Update "Shell Completion" instruction on Linux
   * Update FAQ page
4. misc:
   * log: enable override coloring based on `CLICOLOR` and `CLICOLOR_FORCE`

### v1.45.2

1. misc:
   * fix: help command

### v1.45.1

1. updated linters:
   * `interfacer`: inactivate with go1.18
   * `govet`: inactivate unsupported analyzers (go1.18)
   * `depguard`: reduce requirements
   * `structcheck`: inactivate with go1.18
   * `varnamelen`: bump from v0.6.0 to v0.6.1
2. misc:
   * Automatic Go version detection 🎉 (go1.18)
   * docker: update base images (go1.18)

### v1.45.0

1. updated linters:
   * `cobra`: from 1.3.0 to 1.4.0
   * `containedctx`: from 1.0.1 to 1.0.2
   * `errcheck`: add an option to remove default exclusions
   * `gci`: from 0.3.1 to 0.3.2
   * `go-header`: from 0.4.2 to 0.4.3
   * `gofumpt`: add module-path setting
   * `gofumpt`: from 0.2.1 to 0.3.0
   * `gopsutil`: from 3.22.1 to 3.22.2
   * `gosec`: from 2.9.6 to 2.10.0
   * `makezero`: from 1.1.0 to 1.1.1
   * `revive`: fix default values
   * `wrapcheck`: from 2.4.0 to 2.5.0
2. documentation:
   * docs: add "back to the top" button
   * docs: add `forbidigo` example that uses comments
   * docs: improve linters page
3. misc:
   * go1.18 support 🎉
   * Add an option to manage the targeted version of Go
   * Default to YAML when config file has no extension

### v1.44.2

1. updated linters:
   * `gci`: bump to HEAD
   * `gci`: restore defaults for sections
   * `whitespace`: from 0.0.4 to 0.0.5
2. documentation:
   * add link to configuration in the linters list

### v1.44.1

1. updated linters:
   * `bidichk`: from 0.2.1 to 0.2.2
   * `errchkjson`: from 0.2.1 to 0.2.3
   * `thelper`: from 0.5.0 to 0.5.1
   * `tagliatelle`: from 0.3.0 to 0.3.1
   * `gopsutil`: from 3.21.12 to 3.22.1
   * `gci`: from 0.2.9 to 0.3.0
   * `revive`: from v1.1.3 to v1.1.4
   * `varnamelen`: from v0.5.0 to v0.6.0
2. documentation:
   * linters: improve configuration pages
   * `decorder`: fix `disable-init-func-first-check: false` elaboration
3. misc:
   * fix debug output

### v1.44.0

1. new linters:
   * `containedctx`: https://github.com/sivchari/containedctx
   * `decorder`: https://gitlab.com/bosi/decorder
   * `errchkjson`: https://github.com/breml/errchkjson
   * `maintidx`: https://github.com/yagipy/maintidx
   * `grouper`: https://github.com/leonklingele/grouper
2. updated linters:
   * `asciicheck`: bump to v0.1.1
   * `bidichk`: from 0.1.1 to 0.2.1
   * `bodyclose`: bump to HEAD
   * `decorder`: from 0.2.0 to 0.2.1
   * `depguard`: from 1.0.1 to 1.1.0
   * `errchkjson`: from 0.2.0 to 0.2.1
   * `errorlint`: bump to HEAD
   * `exhaustive`: drop deprecated/unused settings
   * `exhaustive`: from v0.2.3 to 0.7.11
   * `forbidigo`: from 1.2.0 to 1.3.0
   * `forcetypeassert`: bump to v0.1.0
   * `gocritic`: from 0.6.1 to 0.6.2
   * `gocritic`: support autofix
   * `gocyclo`: from 0.3.1 to 0.4.0
   * `godot`: add period option
   * `gofumpt`: from 0.1.1 to 0.2.1
   * `gomnd`: from 2.4.0 to 2.5.0
   * `gomnd`: new configuration
   * `gosec`: from 2.9.1 to 2.9.6
   * `ifshort`: from 1.0.3 to 1.0.4
   * `ineffassign`: bump to HEAD
   * `makezero`: to v1.1.0
   * `promlinter`: from v0.1.0 to HEAD
   * `revive`: fix `enableAllRules`
   * `revive`: from 1.1.2 to 1.1.3
   * `staticcheck`: from 0.2.1 to 0.2.2
   * `tagliatelle`: from 0.2.0 to 0.3.0
   * `thelper`: from 0.4.0 to 0.5.0
   * `unparam`: bump to HEAD
   * `varnamelen`: bump to v0.5.0
   * `wrapcheck`: update configuration to include `ignoreSignRegexps`
3. documentation:
   * linters: improve pages about configuration
   * improve page about false-positive
   * `nolintlint`: fix wrong default value in comment
   * `revive`: add a more detailed configuration
4. misc:
   * outputs: Add support for multiple outputs
   * outputs: Print error text in `<failure>` tag content for more readable JUnit output
   * outputs: ensure that the Issues key in JSON format is a list
   * Return error if any linter fails to run
   * cli: Show deprecated mark in the CLI linters help

### November 2021

1. new linters:
   * `bidichk`: https://github.com/breml/bidichk
2. update linters:
   * `nestif`: from 0.3.0 to 0.3.1
   * `rowserrcheck`: from 1.1.0 to 1.1.1
   * `gopsutil`: from 3.21.9 to 3.21.10
   * `wrapcheck`: from 2.3.1 to 2.4.0
   * `gocritic`: add support for variable substitution in `ruleguard` path settings
3. documentation:
   * improve `go-critic` documentation
   * improve `nolintlint` documentation
4. Misc:
   * cli: don't hide `enable-all` option

### october 2021

1. new linters:
   * `contextcheck`: https://github.com/kkHAIKE/contextcheck
   * `varnamelen`: https://github.com/blizzy78/varnamelen
2. update linters:
   * `gochecknoglobals`: to v0.1.0
   * `gosec`: filter issues according to the severity and confidence
   * `errcheck`: empty selector name.
   * `ifshort`: from 1.0.2 to 1.0.3
   * `go-critic`: from 0.5.6 to 0.6.0
   * `gosec`: from 2.8.1 to 2.9.1
   * `durationcheck`: from 0.0.8 to 0.0.9
   * `wrapcheck`: from 2.3.0 to 2.3.1
   * `revive`: from 1.1.1 to 1.1.2

### September 2021

1. new linters:
   * `ireturn`: https://github.com/butuzov/ireturn
   * `nilnil`: https://github.com/Antonboom/nilnil
   * `tenv`: https://github.com/sivchari/tenv
2. update linters:
   * `errcheck`: update to HEAD
   * `errname`: from 0.1.4 to 0.1.5
   * `gci`: Parse the settings more similarly to the CLI
   * `godot`: from 1.4.9  to 1.4.11
   * `ireturn`: from 0.1.0 to 0.1.1
   * `nlreturn`: add block-size option
   * `paralleltest`: from 1.0.2 to 1.0.3
3. Misc:
   * new-from-rev: add support for finding issues in entire files in a diff

### August 2021

1. new linters:
   * `errname`: https://github.com/Antonboom/errname
2. update linters:
   * `errname`: from 0.1.3 to 0.1.4
   * `go-critic`: fix invalid type conversions.
   * `godot`: from 1.4.8 to 1.4.9
   * `gomodguard`: from 1.2.2 to 1.2.3
   * `revive`: from 1.0.9 to 1.1.1
   * `staticcheck`: bump to 2021.1.1 (v0.2.1)
   * `wrapcheck`: bump to v2.3.0
3. Misc:
   * build binaries and Docker images with go1.17

### July 2021

1. update linters:
   * `errcheck`: allow exclude config without extra file
   * `exhaustive`: from 0.1.0 to 0.2.3
   * `gocognit`: from 1.0.1 to 1.0.5
   * `godot`: from 1.4.7 to 1.4.8
   * `gomoddirectives`: from 0.2.1 to 0.2.2
   * `revive`: from 1.0.8 to 1.0.9
2. documentation:
   * improve `goconst` documentation
   * improve `goimports` description

### June 2021

1. update linters:
   * `durationcheck`: from 0.0.7 to 0.0.8
   * `gci`: from 0.2.8 to 0.2.9
   * `goconst`: from 0.5.6 to 0.5.7
   * `gofumpt`: Add lang-version option
   * `gomodguard`: from 1.2.1 to 1.2.2
   * `gosec`: from 2.8.0 to 2.8.1
   * `revive`: add enable-all-rules.
   * `revive`: allow to disable rule
   * `revive`: fix exclude comment rule for const block
   * `revive`: from 1.0.7 to 1.0.8
   * `wrapcheck`: from 2.1.0 to 2.2.0
2. documentation:
   * add all integrations to docs introduction page
3. Misc:
   * 🎉 Un-deprecate enable-all option
   * output: generate HTML report
   * Support RISV64

### May 2021

1. new linters:
   * `tagliatelle`: https://github.com/ldez/tagliatelle
   * `promlinter`: https://github.com/yeya24/promlinter
2. update linters:
   * `durationcheck`: from 0.0.6 to 0.0.7
   * `errorlint`: bump to HEAD
   * `forbidigo`: from 1.1.0 to 1.2.0
   * `go-critic`: from 0.5.5 to 0.5.6
   * `godot`: from 1.4.6 to 1.4.7
   * ⚠ `golint`: deprecated
   * `gomnd`: from 2.3.2 to 2.4.0
   * `gomodguard`: fix problem where duplicate issues were reported
   * `gosec`: from 2.7.0 to 2.8.0
   * `govet`: fix `sigchanyzer`
   * `govet`: Update vet passes
   * `importas`: allow repeated aliases
   * `importas`: bump to HEAD
   * `makezero`: bump to HEAD
   * `nolintlint`: fix false positive
   * `revive`: convert hard coded excludes into default exclude patterns
   * `revive`: fix add-constant rule support
   * `revive`: fix excludes
   * `revive`: from 1.0.6 to 1.0.7
   * `revive`: improve 'exported' rule output
   * `rowserrcheck`: bump to v1.1.0
   * `staticcheck`: configuration for `staticcheck`, `gosimple`, `stylecheck`
   * `staticcheck`: from 0.1.3 to 0.1.4
   * `staticcheck`: from v0.1.4 to v0.2.0
   * `wastedassign`: from 0.2.0 to 1.0.0
   * `wastedassign`: from 1.0.0 to v2.0.6
   * `wrapcheck`: from 1.2.0 to 2.1.0
3. documentation:
   * improve linters page
   * `exhaustivestruct` example explanation
   * fix pattern of `forbidigo` in example config yaml
   * bump documentation dependencies
   * fix typos
4. Misc:
   * set the minimum Go version to go1.15
   * non-zero exit code when a linter produces a panic

### April 2021

1. new linters:
   * `tagliatelle`: https://github.com/ldez/tagliatelle
   * `promlinter`: https://github.com/yeya24/promlinter
2. update linters:
   * `godot`: from 1.4.4 to 1.4.6
   * `wrapcheck`: from 1.0.0 to 1.2.0
   * `go-mnd`: from 2.3.1 to 2.3.2
   * `wsl`: from 3.2.0 to 3.3.0
   * `revive`: from 1.0.5 to 1.0.6
   * `importas`: bump to HEAD
   * `staticcheck`: configurable Go version
   * `gosec`: add configuration
   * `typecheck`: improve error stack parsing
3. documentation:
   * bump documentation dependencies
   * fix typos
4. Misc:
   * fix: comma in exclude pattern leads to unexpected results

### March 2021

1. new linters:
   * `gomoddirectives`: https://github.com/ldez/gomoddirectives
2. update linters:
   * `go-critic`: from 0.5.4 to 0.5.5
   * `gofumpt`: from v0.1.0 to v0.1.1
   * `gosec`: from 2.6.1 to 2.7.0
   * `ifshort`: bump to v1.0.2
   * `importas`: bump to HEAD
   * `makezero`: bump to HEAD
   * `nolintlint`: allow to fix //nolint lines
   * `revive`: from 1.0.3 to 1.0.5
   * `revive`: the default configuration is only applied when no dedicated configuration
   * `rowserrcheck`: bump to HEAD
   * ⚠ `scopelint`: deprecated
   * `staticcheck`: from v0.1.2 (2020.2.2) to v0.1.3 (2020.2.3)
   * 🎉 `typecheck`: display compilation errors as report instead of error
   * `wastedassign`: bump to v0.2.0
   * `wrapcheck`: bump to v1.0.0
3. documentation:
   * improve [linters page](https://golangci-lint.run/usage/linters/) (versions, deprecation, and presets)
   * add [cache directory](https://golangci-lint.run/usage/configuration/#cache) information
   * adding missing format options
   * fix typos
4. Misc:
   * Set `version` command output to Stdout
   * fix linters load mode
   * Restore fast linters meaning

### February 2021

1. new linters:
   * `durationcheck`: https://github.com/charithe/durationcheck
   * `revive`: https://github.com/mgechev/revive
   * `cyclop`: https://github.com/bkielbasa/cyclop
   * `wastedassign`: https://github.com/sanposhiho/wastedassign
   * `importas`: https://github.com/julz/importas
   * `nilerr`: https://github.com/gostaticanalysis/nilerr
   * `forcetypeassert`: https://github.com/gostaticanalysis/forcetypeassert
   * `govet`: add `stringintconv` and `ifaceassert`
2. update linters:
   * `prealloc`: Use upstream version
   * `errcheck`: Use upstream version
   * `ineffassign`: Use upstream version
   * `gocyclo`: Use upstream version
   * `godot` from 1.3.2 to 1.4.3
   * `exhaustivestruct` from 1.1.0 to 1.2.0
   * `forbidigo` from 1.0.0 to 1.1.0
   * `thelper` from 0.2.1 to 0.3.1
   * `ruleguard`: print error message and exit with non-zero status when parsing error occurs
   * fix wrong load mode for `asciicheck`, `exhaustivestruct`, `exportloopref`, and `makezero`
   * `wsl`: bump to v3.2.0
   * `durationcheck`: from 0.0.4 to 0.0.6
   * `staticcheck`: from 2020.1.6 to v0.1.2 (2020.2.2)
   * `thelper` from 0.3.1 to 0.4.0
   * `bodyclose`: bump to HEAD
   * `go-err113`: bump to HEAD
   * ⚠ `interfacer`: deprecated
   * ⚠ `maligned`: deprecated (govet `fieldalignment` as replacement)
   * `govet`: use same default linter as go vet
   * `go-printf-func-name`: to `v0.0.0-20200119135958-7558a9eaa5af`
   * `godox`: to `v0.0.0-20210227103229-6504466cf951`
   * `asciicheck`: to `v0.0.0-20200416200610-e657995f937b`
   * `wrapcheck`: to `v0.0.0-20201130113247-1683564d9756`
   * `unparam`: to `v0.0.0-20210104141923-aac4ce9116a7`
3. CLI: truncate multiline descriptions
4. fix: new-from-rev for a large repository
5. Support RelatedInformation for analysis Diagnostic
6. use go1.16 to create binaries
7. fix: MIPS release
8. documentation:
   * bump documentation dependencies
   * fix `go-header` usage
   * improve `gocritic` description
   * update deprecated hyperlink for Sublime Text plugin
   * add docs on using homebrew tap

### January 2021

1. new linters:
   * `predeclared`: https://github.com/nishanths/predeclared
   * `ifshort`: https://github.com/esimonov/ifshort
2. update linters:
   * `go-critic` from 0.5.2 to 0.5.3
   * `thelper` from 0.1.0 to 0.2.1
   * Validate `go-critic` settings
   * `gofumpt` to v0.1.0
   * `gci` to v0.2.8
   * `go-mnd` to v2.3.1
   * `gosec` from 2.5.0 to 2.6.1
   * `godot` from 1.3.2 to 1.4.3
   * `ifshort` to v1.0.1
   * `rowserrcheck`: fix reports false positive
3. fix: modules-download-mode support
4. documentation:
   * bump documentation dependencies

### December 2020

1. new linters:
   * `forbidigo`: https://github.com/ashanbrown/forbidigo
   * `makezero`: https://github.com/ashanbrown/makezero
   * `thelper`: https://github.com/kulti/thelper
2. update linters:
   * `go-header` from v0.3.1 to v0.4.2
   * `go-mnd` from v2.0.0 to v2.2.0
   * `godot` from v1.3.0 to v1.3.2
   * `gci` from v0.2.4 to v0.2.7
   * `gomodguard` from v1.1.0 to v1.2.0
   * `go-errorlint` from v0.0.0-20201006195004-351e25ade6e3 to v0.0.0-20201127212506-19bd8db6546f
   * `gofumpt` from v0.0.0-20200802201014-ab5a8192947d to v0.0.0-20201129102820-5c11c50e9475
   * `nolintlint` fix comment analysis. (#1571)
3. result/processors: treat all non-Go source as special autogenerated files
4. throw an error on panic. (#1540)
5. resolve custom linters' path relative to config file directory (#1572)
6. treat all non-Go source as special autogenerated files
7. documentation:
   * add settings examples for `gocritic` (#1562)
   * removing reference to no-longer-existing linter-in-the-cloud (#1553)
8. others:
   * bump `gopkg.in/yaml.v2` from 2.3.0 to 2.4.0 (#1528)
   * bump `gatsby-remark-responsive-iframe` in /docs (#1533)
   * bump `gatsby-remark-images` from 3.3.29 to 3.6.0 in /docs (#1531)
   * bump `ini` from 1.3.5 to 1.3.8 in /tools (#1560)
   * bump `react-headroom` from 3.0.0 to 3.0.1 in /docs (#1532)
   * bump `react-live` from 2.2.2 to 2.2.3 in /docs (#1534)
   * bump `react` from 16.13.1 to 16.14.0 in /docs (#1481)
   * Fix `forbidigo` linter name in reports (#1590)

### November 2020

1. new linters:
   * `paralleltest`: https://github.com/kunwardeep/paralleltest
2. update linters:
   * `godot` from v0.4.9 to v1.3.0
   * `exportloopref` from v0.1.7 to v0.1.8
   * `gosec` from 2.4.0 to 2.5.0
   * `goconst` using upstream https://github.com/jgautheron/goconst
3. `DefaultExcludePatterns` should only be used for specified linter (#1494)
4. unknown linter breaks //nolint (#1497)
5. report all unknown linters at once (#1477)
6. CI:
   * fix Docker tag for Alpine build
7. documentation:
   * missing sort-results in the docs (#1514)
   * add description of Homebrew's official formula (#1421)
8. others:
   * bump `golang.org/x/text` to v0.3.4 (#1293)
   * bump `github.com/fatih/color` to from 1.9.0 to 1.10.0 (#1485)
   * bump `lodash` from 4.17.15 to 4.17.19 in /.github/peril (#1252)
   * bump `polished` from 3.6.6 to 4.0.3 in /docs (#1482)
   * bump `gatsby-alias-imports` from 1.0.4 to 1.0.6 in /docs (#1479)
   * bump `puppeteer` from 5.3.1 to 5.4.1 in /docs (#1480)
   * bump `gatsby-remark-embedder` from 3.0.0 to 4.0.0 in /docs (#1478)

### October 2020

1. new linters:
   * `exhaustivestruct`: https://github.com/mbilski/exhaustivestruct
   * `go-errorlint`: https://github.com/polyfloyd/go-errorlint
   * `tparallel`: https://github.com/moricho/tparallel
   * `wrapcheck`: https://github.com/tomarrell/wrapcheck
2. update linters:
   * `honnef.co/go/tools` from 2020.1.5 to 2020.1.6
   * `exhaustivestruct` from v1.0.1 to v1.1.0
   * `exhaustive` to v0.1.0
   * `gochecknoglobals`: use https://github.com/leighmcculloch/gochecknoglobals
3. add support for powershell completion (#1408)
4. add `.golangci.yaml` to list of configuration files searched on startup (#1364)
5. support for only specifying default severity (#1396)
6. documentation:
   * mention macports installation procedure on macOS (#1352)
   * sort linters (#1451)
7. CI:
   * add codeQL scanning (#1405)
   * fix version details in Docker image (#1471)
   * releasing docker image for arm64 (#1383)
   * change interval for npm to monthly (#1424)
8. others:
   * use tag version for cobra (#1458)
   * bump `nancy` to 1.0.1 (#1410)
   * bump `gatsby-plugin-catch-links` in /docs (#1415)
   * bump `gatsby-plugin-mdx` from 1.2.40 to 1.2.43 in /docs (#1419)
   * bump `gatsby-plugin-sharp` from 2.6.31 to 2.6.40 in /docs (#1423)
   * bump `gatsby-plugin-sitemap` from 2.4.12 to 2.4.14 in /docs (#1417)
   * bump `github.com/mattn/go-colorable` from 0.1.7 to 0.1.8 (#1413)
   * bump `github.com/sirupsen/logrus` from 1.6.0 to 1.7.0 (#1412)
   * bump `github.com/sourcegraph/go-diff` from 0.6.0 to 0.6.1 (#1414)
   * bump `golangci/golangci-lint-action` from v2 to v2.3.0  (#1447) (#1469)
   * bump `puppeteer` from 3.3.0 to 5.3.1 in /docs (#1418)

### September 2020

1. update linters:
   * `godot` from 0.4.8 to 0.4.9
   * `exhaustive` from v0.0.0-20200708172631-8866003e3856 to v0.0.0-20200811152831-6cf413ae40e0
   * `gofumpt` from v0.0.0-20200709182408-4fd085cb6d5f to v0.0.0-20200802201014-ab5a8192947d
2. add support for fish completion (#1201)
3. documentation:
   * fix typo in performance docs (#1350)
4. CI:
   * prevent macos to be marked as passing upon failure (#1381)
   * check only for go.mod file (#1397)
   * check if go.mod and go.sum are up to dated (#1377)
   * trigger Netlify (#1358)
5. others:
   * bump `github.com/sourcegraph/go-diff` from 0.5.3 to 0.6.0 (#1353)
   * bump `github.com/valyala/quicktemplate` from 1.6.2 to 1.6.3 (#1385)
   * ignore known dependency failure in nancy (#1378)
   * bump `@mdx-js/mdx` from 1.6.16 to 1.6.18 in /docs (#1401)
   * bump `gatsby` from 2.24.52 to 2.24.65 in /docs (#1400)
   * bump `gatsby-plugin-canonical-urls` in /docs (#1390)
   * bump `gatsby-plugin-google-analytics` in /docs (#1388)
   * bump `gatsby-plugin-manifest` from 2.4.23 to 2.4.27 in /docs (#1355)
   * bump `gatsby-plugin-mdx` from 1.2.35 to 1.2.40 in /docs (#1386)
   * bump `gatsby-plugin-offline` from 3.2.23 to 3.2.27 in /docs (#1368)
   * bump `gatsby-plugin-sharp` from 2.6.25 to 2.6.31 in /docs (#1354)
   * bump `gatsby-plugin-sitemap` from 2.4.11 to 2.4.12 in /docs (#1344)
   * bump `gatsby-remark-autolink-headers` in /docs (#1387)
   * bump `gatsby-remark-images` from 3.3.25 to 3.3.28 in /docs (#1345)
   * bump `gatsby-remark-images` from 3.3.28 to 3.3.29 in /docs (#1365)
   * bump `gatsby-remark-mermaid` from 2.0.0 to 2.1.0 in /docs (#1369)
   * bump `gatsby-source-filesystem` in /docs (#1366)
   * bump `gatsby-source-filesystem` in /docs (#1389)
   * bump `gatsby-transformer-sharp` in /docs (#1402)
   * bump `gatsby-transformer-yaml` from 2.4.10 to 2.4.11 in /docs (#1367)
   * bump `node-fetch` in /.github/contributors (#1363)
   * bump `polished` from 3.6.5 to 3.6.6 in /docs (#1347)

### August 2020

1. new `nlreturn` linter: https://github.com/ssgreg/nlreturn
2. new `gci` linter: https://github.com/daixiang0/gci
3. support `latest` version of golangci-lint in golangci-lint-action
4. update `gosec` linter from 2.3.0 to 2.4.0
5. update `godot` linter from 0.4.2 to 0.4.8
6. update `go-critic` from 0.5.0 to 0.5.2 (#1307)
7. update `nlreturn` from 2.0.1 to 2.0.2 (#1287), 2.0.2 to 2.1.0 (#1327)
8. update `gci` to v0.2.1 (#1292), to v0.2.2 (#1305), to v0.2.4 (#1337),
9. update `funlen` from 0.0.2 to 0.0.3 (#1341)
10. upgrade to golang 1.15 for smaller binary (#1303)
11. support short and json formats for version cmd (#1315)
12. add home directory to config file search paths (#1325)
13. allow for serializing multiple golangci-lint invocations (#1302)

### July 2020

1. `gofumpt` linter:
    * update linter
    * add `extra-rules` option
    * support auto-fixing
2. upgrade `exhaustive` linter
3. upgrade `exportloopref` linter
4. improve 'no such linter' error message
5. sorting result.Issues implementation
6. enhancements in CI:
    * Run `nancy` validation for all dependencies
    * Move dependabot config to `.github` folder
7. other
    * bump `lodash` from 4.17.15 to 4.17.19 in /tools
    * bump `golangci/golangci-lint-action` from v1.2.2 to v2
    * bump `github.com/valyala/quicktemplate` from 1.5.0 to 1.5.1


### June 2020

1. Add new linters: `gofumpt`

### May 2020

1. Add new linters: `nolintlint`, `goerr113`
2. Updated linters: `godot`, `staticcheck`
3. Launch a [website](https://golangci-lint.run)

### April 2020

1. Add new linters: `testpackage`, `nestif`, `godot`, `gomodguard`, `asciicheck`
2. Add GitHub Actions output format
3. Update linters: `wsl`, `gomodguard`, `gosec`
4. Support `disabled-tags` setting for `gocritic`
5. Mitigate OOM and "failed prerequisites"
6. Self-isolate due to unexpected pandemics
7. Support case-sensitive excludes
8. Allow granular re-enabling excludes by ID, e.g. `EXC0002`

### September 2019

1. Support go1.13
2. Add new linters: `funlen`, `whitespace` (with auto-fix) and `godox`
3. Update linters: `gochecknoglobals`, `scopelint`, `gosec`
4. Provide pre-built binary for ARM and FreeBSD
5. Fix false-positives in `unused`
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
