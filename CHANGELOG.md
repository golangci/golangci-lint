Follow the news and releases on [Mastodon](https://fosstodon.org/@golangcilint) and on [Bluesky](https://bsky.app/profile/golangci-lint.run).

`golangci-lint` is a free and open-source project built by volunteers.

If you value it, consider supporting us, we appreciate it!

[![Donate](https://img.shields.io/badge/Donate-❤️-blue?style=for-the-badge)](https://donate.golangci.org)

<!-- START --->

### v2.7.2

_Released on 2025-12-07_

1. Linter bug fixes
   * `gosec`: from 2.22.10 to daccba6b93d7

### v2.7.1

_Released on 2025-12-04_

1. Linter bug fixes
   * `modernize`: disable `stringscut` analyzer

### v2.7.0

_Released on 2025-12-03_

1. Bug fixes
   * fix: clone args used by `custom` command
2. Linters new features or changes
   * `no-sprintf-host-port`: from 0.2.0 to 0.3.1 (ignore string literals without a colon)
   * `unqueryvet`: from 1.2.1 to 1.3.0 (handles `const` and `var` declarations)
   * `revive`: from 1.12.0 to 1.13.0 (new option: `enable-default-rules`, new rules: `forbidden-call-in-wg-go`, `unnecessary-if`, `inefficient-map-lookup`)
   * `modernize`: from 0.38.0 to 0.39.0 (new analyzers: `plusbuild`, `stringscut`)
3. Linters bug fixes
   * `perfsprint`: from 0.10.0 to 0.10.1
   * `wrapcheck`: from 2.11.0 to 2.12.0
   * `godoc-lint`: from 0.10.1 to 0.10.2
4. Misc.
   * Add some flags to the `custom` command
5. Documentation
   * docs: split changelog v1 and v2

### v2.6.2

_Released on 2025-11-14_

1. Bug fixes
   * `fmt` command with symlinks
   * use file depending on build configuration to invalidate cache
2. Linters bug fixes
   * `testableexamples`: from 1.0.0 to 1.0.1
   * `testpackage`: from 1.1.1 to 1.1.2

### v2.6.1

_Released on 2025-11-04_

1. Linters bug fixes
   * `copyloopvar`: from 1.2.1 to 1.2.2
   * `go-critic`: from 0.14.0 to 0.14.2

### v2.6.0

_Released on 2025-10-29_

1. New linters
   * Add `modernize` analyzer suite
2. Linters new features or changes
   * `arangolint`: from 0.2.0 to 0.3.1
   * `dupword`: from 0.1.6 to 0.1.7 (new option `comments-only`)
   * `go-critic`: from 0.13.0 to 0.14.0 (new rules/checkers: `zeroByteRepeat`, `dupOption`)
   * `gofumpt`: from 0.9.1 to 0.9.2 ("clothe" naked returns is now controlled by the `extra-rules` option)
   * `perfsprint`: from 0.9.1 to 0.10.0 (new options: `concat-loop`, `loop-other-ops`)
   * `wsl`: from 5.2.0 to 5.3.0
3. Linters bug fixes
   * `dupword`: from 0.1.6 to 0.1.7
   * `durationcheck`: from 0.0.10 to 0.0.11
   * `exptostd`: from 0.4.4 to 0.4.5
   * `fatcontext`: from 0.8.1 to 0.9.0
   * `forbidigo`: from 2.1.0 to 2.3.0
   * `ginkgolinter`: from 0.21.0 to 0.21.2
   * `godoc-lint`: from 0.10.0 to 0.10.1
   * `gomoddirectives`: from 0.7.0 to 0.7.1
   * `gosec`: from 2.22.8 to 2.22.10
   * `makezero`: from 2.0.1 to 2.1.0
   * `nilerr`: from 0.1.1 to 0.1.2
   * `paralleltest`: from 1.0.14 to 1.0.15
   * `protogetter`: from 0.3.16 to 0.3.17
   * `unparam`: from 0df0534333a4 to 5beb8c8f8f15
4. Misc.
   * fix: ignore some files to hash the version for custom build

### v2.5.0

_Released on 2025-09-21_

1. New linters
   * Add `godoclint` linter https://github.com/godoc-lint/godoc-lint
   * Add `unqueryvet` linter https://github.com/MirrexOne/unqueryvet
   * Add `iotamixing` linter https://github.com/AdminBenni/iota-mixing
2. Linters new features or changes
   * `embeddedstructfieldcheck`: from 0.3.0 to 0.4.0 (new option: `empty-line`)
   * `err113`: from aea10b59be24 to 0.1.1 (skip internals of `Is` methods for `error` type)
   * `ginkgolinter`: from 0.20.0 to 0.21.0 (new option: `force-tonot`)
   * `gofumpt`: from 0.8.0 to 0.9.1 (new rule is to "clothe" naked returns for the sake of clarity)
   * `ineffassign`: from 0.1.0 to 0.2.0 (new option: `check-escaping-errors`)
   * `musttag`: from 0.13.1 to 0.14.0 (support interface methods)
   * `revive`: from 1.11.0 to 1.12.0 (new options: `identical-ifelseif-branches`, `identical-ifelseif-conditions`, `identical-switch-branches`, `identical-switch-conditions`, `package-directory-mismatch`, `unsecure-url-scheme`, `use-waitgroup-go`, `useless-fallthrough`)
   * `thelper`: from 0.6.3 to 0.7.1 (skip `t.Helper` in functions passed to `synctest.Test`)
   * `wsl`: from 5.1.1 to 5.2.0 (improvements related to subexpressions)
3. Linters bug fixes
   * `asciicheck`: from 0.4.1 to 0.5.0
   * `errname`: from 1.1.0 to 1.1.1
   * `fatcontext`: from 0.8.0 to 0.8.1
   * `go-printf-func-name`: from 0.1.0 to 0.1.1
   * `godot`: from 1.5.1 to 1.5.4
   * `gosec`: from 2.22.7 to 2.22.8
   * `nilerr`: from 0.1.1 to a temporary fork
   * `nilnil`: from 1.1.0 to 1.1.1
   * `protogetter`: from 0.3.15 to 0.3.16
   * `tagliatelle`: from 0.7.1 to 0.7.2
   * `testifylint`: from 1.6.1 to 1.6.4
4. Misc.
   * fix: "no export data" errors are now handled as a standard typecheck error
5. Documentation
   * Improve nolint section about syntax

### v2.4.0

_Released on 2025-08-14_

1. Enhancements
    * 🎉 go1.25 support
2. Linters new features or changes
    * `exhaustruct`: from v3.3.1 to 4.0.0 (new options: `allow-empty`, `allow-empty-rx`, `allow-empty-returns`, `allow-empty-declarations`)
3. Linters bug fixes
   * `godox`: trim filepath from report messages
   * `staticcheck`: allow empty options
   * `tagalign`: from 1.4.2 to 1.4.3
4. Documentation
   * 🌟 New website (with a search engine)

### v2.3.1

_Released on 2025-08-02_

1. Linters bug fixes
   * `gci`: from 0.13.6 to 0.13.7
   * `gosec`: from 2.22.6 to 2.22.7
   * `noctx`: from 0.3.5 to 0.4.0
   * `wsl`: from 5.1.0 to 5.1.1
   * tagliatelle: force upper case for custom initialisms

### v2.3.0

_Released on 2025-07-21_

1. Linters new features or changes
   * `ginkgolinter`: from 0.19.1 to 0.20.0 (new option: `force-assertion-description`)
   * `iface`: from 1.4.0 to 1.4.1 (report message improvements)
   * `noctx`: from 0.3.4 to 0.3.5 (new detections: `log/slog`, `exec`, `crypto/tls`)
   * `revive`: from 1.10.0 to 1.11.0 (new rule: `enforce-switch-style`)
   * `wsl`: from 5.0.0 to 5.1.0
2. Linters bug fixes
   * `gosec`: from 2.22.5 to 2.22.6
   * `noinlineerr`: from 1.0.4 to 1.0.5
   * `sloglint`: from 0.11.0 to 0.11.1
3. Misc.
   * fix: panic close of closed channel

### v2.2.2

_Released on 2025-07-11_

1. Linters bug fixes
   * `noinlineerr`: from 1.0.3 to 1.0.4
2. Documentation
   * Improve debug keys documentation
3. Misc.
   * fix: panic close of closed channel
   * godot: add noinline value into the JSONSchema

### v2.2.1

_Released on 2025-06-28_

1. Linters bug fixes
  * `varnamelen`: fix configuration

### v2.2.0

_Released on 2025-06-28_

1. New linters
   * Add `arangolint` linter https://github.com/Crocmagnon/arangolint
   * Add `embeddedstructfieldcheck` linter https://github.com/manuelarte/embeddedstructfieldcheck
   * Add `noinlineerr` linter https://github.com/AlwxSin/noinlineerr
   * Add `swaggo` formatter https://github.com/golangci/swaggoswag
2. Linters new features or changes
   * `errcheck`: add `verbose` option
   * `funcorder`: from 0.2.1 to 0.5.0 (new option `alphabetical`)
   * `gomoddirectives`: from 0.6.1 to 0.7.0 (new option `ignore-forbidden`)
   * `iface`: from 1.3.1 to 1.4.0 (new option `unexported`)
   * `noctx`: from 0.1.0 to 0.3.3 (new report messages, and new rules related to `database/sql`)
   * `noctx`: from 0.3.3 to 0.3.4 (new SQL functions detection)
   * `revive`: from 1.9.0 to 1.10.0 (new rules: `time-date`, `unnecessary-format`, `use-fmt-print`)
   * `usestdlibvars`: from 1.28.0 to 1.29.0 (new option `time-date-month`)
   * `wsl`: deprecation
   * `wsl_v5`: from 4.7.0 to 5.0.0 (major version with new configuration)
3. Linters bug fixes
   * `dupword`: from 0.1.3 to 0.1.6
   * `exptostd`: from 0.4.3 to 0.4.4
   * `forbidigo`: from 1.6.0 to 2.1.0
   * `gci`: consistently format the code
   * `go-spancheck`: from 0.6.4 to 0.6.5
   * `goconst`: from 1.8.1 to 1.8.2
   * `gosec`: from 2.22.3 to 2.22.4
   * `gosec`: from 2.22.4 to 2.22.5
   * `makezero`: from 1.2.0 to 2.0.1
   * `misspell`: from 0.6.0 to 0.7.0
   * `usetesting`: from 0.4.3 to 0.5.0
4. Misc.
   * exclusions:  fix `path-expect`
   * formatters: write the input to `stdout` when using `stdin` and there are no changes
   * migration: improve the error message when trying to migrate a migrated config
   * `typecheck`: deduplicate errors
   * `typecheck`: stops the analysis after the first error
   * Deprecate `print-resources-usage` flag
   * Unique version per custom build
5. Documentation
   * Improves typecheck FAQ
   * Adds plugin systems recommendations
   * Add description for `linters.default` sets

### v2.1.6

_Released on 2025-05-04_

1. Linters bug fixes
   * `godot`: from 1.5.0 to 1.5.1
   * `musttag`: from 0.13.0 to 0.13.1
2. Documentation
   * Add note about golangci-lint v2 integration in VS Code

### v2.1.5

_Released on 2025-04-24_

Due to an error related to Snapcraft, some artifacts of the v2.1.4 release have not been published.

This release contains the same things as v2.1.3.

### v2.1.4

_Released on 2025-04-24_

Due to an error related to Snapcraft, some artifacts of the v2.1.3 release have not been published.

This release contains the same things as v2.1.3.

### v2.1.3

_Released on 2025-04-24_

1. Linters bug fixes
   * `fatcontext`: from 0.7.2 to 0.8.0
2. Misc.
   * migration: fix `nakedret.max-func-lines: 0`
   * migration: fix order of `staticcheck` settings
   * fix: add `go.mod` hash to the cache salt
   * fix: use diagnostic position for related information position

### v2.1.2

_Released on 2025-04-15_

1. Linters bug fixes
   * `exptostd`: from 0.4.2 to 0.4.3
   * `gofumpt`: from 0.7.0 to 0.8.0
   * `protogetter`: from 0.3.13 to 0.3.15
   * `usetesting`: from 0.4.2 to 0.4.3

### v2.1.1

_Released on 2025-04-12_

The release process of v2.1.0 failed due to a regression inside goreleaser.

The binaries of v2.1.0 have been published, but not the other artifacts (AUR, Docker, etc.).

### v2.1.0

_Released on 2025-04-12_

1. Enhancements
   * Add an option to display absolute paths (`--path-mode=abs`)
   * Add configuration path placeholder (`${config-path}`)
   * Add `warn-unused` option for `fmt` command
   * Colored diff for `fmt` command (`golangci-lint fmt --diff-colored`)
2. New linters
   * Add `funcorder` linter https://github.com/manuelarte/funcorder
3. Linters new features or changes
   * `go-errorlint`: from 1.7.1 to 1.8.0 (automatic error comparison and type assertion fixes)
   * ⚠️ `goconst`: `ignore-strings` is deprecated and replaced by `ignore-string-values`
   * `goconst`: from 1.7.1 to 1.8.1 (new options: `find-duplicates`, `eval-const-expressions`)
   * `govet`: add `httpmux` analyzer
   * `nilnesserr`: from 0.1.2 to 0.2.0 (detect more cases)
   * `paralleltest`: from 1.0.10 to 1.0.14 (checks only `_test.go` files)
   * `revive`: from 1.7.0 to 1.9.0 (support kebab case for setting names)
   * `sloglint`: from 0.9.0 to 0.11.0 (autofix, new option `msg-style`, suggest `slog.DiscardHandler`)
   * `wrapcheck`: from 2.10.0 to 2.11.0 (new option `report-internal-errors`)
   * `wsl`: from 4.6.0 to 4.7.0 (cgo files are always excluded)
4. Linters bug fixes
   * `fatcontext`: from 0.7.1 to 0.7.2
   * `gocritic`: fix `importshadow` checker
   * `gosec`: from 2.22.2 to 2.22.3
   * `ireturn`: from 0.3.1 to 0.4.0
   * `loggercheck`: from 0.10.1 to 0.11.0
   * `nakedret`: from 2.0.5 to 2.0.6
   * `nonamedreturns`: from 1.0.5 to 1.0.6
   * `protogetter`: from 0.3.12 to 0.3.13
   * `testifylint`: from 1.6.0 to 1.6.1
   * `unconvert`: update to HEAD
5. Misc.
   * Fixes memory leaks when using go1.(N) with golangci-lint built with go1.(N-X)
   * Adds `golangci-lint-fmt` pre-commit hook
6. Documentation
   * Improvements
   * Updates section about vscode integration

### v2.0.2

_Released on 2025-03-25_

1. Misc.
   * Fixes flags parsing for formatters
   * Fixes the filepath used by the exclusion `source` option
2. Documentation
   * Adds a section about flags migration
   * Cleaning pages with v1 options

### v2.0.1

_Released on 2025-03-24_

1. Linters/formatters bug fixes
   * `golines`: fix settings during linter load
2. Misc.
   * Validates the `version` field before the configuration
   * `forbidigo`: fix migration

### v2.0.0

_Released on 2025-03-24_

1. Enhancements
   * 🌟 New `golangci-lint fmt` command with dedicated [formatter configuration](https://golangci-lint.run/docs/welcome/quick-start/#formatting)
   * ♻️ New `golangci-lint migrate` command to help migration from v1 to v2 (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#command-migrate))
   * ⚠️ New default values (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/))
   * ⚠️ No exclusions by default (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#issuesexclude-use-default))
   * ⚠️ New default sort order (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#outputsort-order))
   * 🌟 New option `run.relative-path-mode` (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#runrelative-path-mode))
   * 🌟 New linters configuration (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#linters))
   * 🌟 New output format configuration (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#output))
   * 🌟 New `--fast-only` flag (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#lintersfast))
   * 🌟 New option `linters.exclusions.warn-unused` to log a warning if an exclusion rule is unused.
2. New linters/formatters
   * Add `golines` formatter https://github.com/segmentio/golines
3. Linters new features
   * ⚠️ Merge `staticcheck`, `stylecheck`, `gosimple` into one linter (`staticcheck`) (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#lintersenablestylecheckgosimplestaticcheck))
   * `go-critic`: from 0.12.0 to 0.13.0
   * `gomodguard`: from 1.3.5 to 1.4.1 (block explicit indirect dependencies)
   * `nilnil`: from 1.0.1 to 1.1.0 (new option: `only-two`)
   * `perfsprint`: from 0.8.2 to 0.9.1 (checker name in the diagnostic message)
   * `staticcheck`: new `quickfix` set of rules
   * `testifylint`: from 1.5.2 to 1.6.0 (new options: `equal-values`, `suite-method-signature`, `require-string-msg`)
   * `wsl`: from 4.5.0 to 4.6.0 (new option: `allow-cuddle-used-in-block`)
4. Linters bug fixes
   * `bidichk`: from 0.3.2 to 0.3.3
   * `errchkjson`: from 0.4.0 to 0.4.1
   * `errname`: from 1.0.0 to 1.1.0
   * `funlen`: fix `ignore-comments` option
   * `gci`: from 0.13.5 to 0.13.6
   * `gosmopolitan`: from 1.2.2 to 1.3.0
   * `inamedparam`: from 0.1.3 to 0.2.0
   * `intrange`: from 0.3.0 to 0.3.1
   * `protogetter`: from 0.3.9 to 0.3.12
   * `unparam`: from 8a5130ca722f to 0df0534333a4
5. Misc.
   * 🧹 Configuration options renaming (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/))
   * 🧹 Remove options (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/))
   * 🧹 Remove flags (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/))
   * 🧹 Remove alternative names (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/#alternative-linter-names))
   * 🧹 Remove or replace deprecated elements (cf. [Migration guide](https://golangci-lint.run/docs/product/migration-guide/))
   * Adds an option to display some commands as JSON:
     * `golangci-lint config path --json`
     * `golangci-lint help linters --json`
     * `golangci-lint help formatters --json`
     * `golangci-lint linters --json`
     * `golangci-lint formatters --json`
     * `golangci-lint version --json`
6. Documentation
   * [Migration guide](https://golangci-lint.run/docs/product/migration-guide/)
