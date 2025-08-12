---
title: Migration guide
weight: 3
aliases:
  - /product/migration-guide/
---

## Command `migrate`

You can use golangci-lint to migrate your configuration with the `migrate` command:

```bash
golangci-lint migrate
```

Be aware that **comments inside a configuration file are not migrated.** You need to add them manually after the migration.

**Deprecated options from v1 or unknown fields are not migrated.**

The migration file format is based on the extension of the [configuration file](/docs/configuration/file).
The format can be overridden by using the `--format` flag:

```bash
golangci-lint migrate --format json
```

Before the migration, the previous configuration file is copied and saved to a file named `<config_file_name>.bck.<config_file_extension>`.

By default, before the migration process, the configuration file is validated against the JSON Schema of configuration v1.
If you want to skip this validation, you can use the `--skip-validation` flag:

```bash
golangci-lint migrate --skip-validation
```

The `migrate` command enforces the following default values:

- `run.timeout`: the existing value is ignored because, in v2, there is no timeout by default.
- `issues.show-stats`: the existing value is ignored because, in v2, stats are enabled by default.
- `run.concurrency`: if the existing value was `0`, it is removed as `0` is the new default.
- `run.relative-path-mode`: if the existing value was `cfg`, it is removed as `cfg` is the new default.

`issues.exclude-generated` has a new default value (v1 `lax`, v2 `strict`), so this field will be added during the migration to maintain the previous behavior.

`issues.exclude-dirs-use-default` has been removed, so it is converted to `linters.exclusions.paths` and, if needed, `formatters.exclusions.paths`.

Other fields explicitly defined in the configuration file are migrated even if the value is the same as the default value.

The `migrate` command automatically migrates `linters.presets` in individual linters to `linters.enable`.

{{% cli-output cmd="migrate" %}}

## Changes

### `linters`

#### `linters.disable-all`

This property has been replaced with `linters.default: none`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  disable-all: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  default: none
```
{{< /tab >}}
{{< /tabs >}}

#### `linters.enable-all`

This property has been replaced with `linters.default: all`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  enable-all: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  default: all
```
{{< /tab >}}
{{< /tabs >}}

#### `linters.enable[].<formatter_name>`

The linters `gci`, `gofmt`, `gofumpt`, and `goimports` have been moved to the `formatters` section.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  enable:
    - gci
    - gofmt
    - gofumpt
    - goimports
```
{{< /tab >}}
{{< tab >}}
```yaml
formatters:
  enable:
    - gci
    - gofmt
    - gofumpt
    - goimports
```
{{< /tab >}}
{{< /tabs >}}

#### `linters.enable[].{stylecheck,gosimple,staticcheck}`

The linters `stylecheck`, `gosimple`, and `staticcheck` has been merged inside the `staticcheck`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  enable:
    - gosimple
    - staticcheck
    - stylecheck
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  enable:
    - staticcheck
```
{{< /tab >}}
{{< /tabs >}}

#### `linters.fast`

This property has been removed.

There are 2 new options (they are not strictly equivalent to the previous option):

1. `linters.default: fast`: set all "fast" linters as the default set of linters.
    ```yaml
    linters:
      default: fast
    ```
2. `--fast-only`: filters all enabled linters to keep only "fast" linters.

#### `linters.presets`

This property has been removed.

The `migrate` command automatically migrates `linters.presets` in individual linters to `linters.enable`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
Presets:

| name | linters |
|------|---------|
| bugs | `asasalint`, `asciicheck`, `bidichk`, `bodyclose`, `contextcheck`, `durationcheck`, `errcheck`, `errchkjson`, `errorlint`, `exhaustive`, `gocheckcompilerdirectives`, `gochecksumtype`, `gosec`, `gosmopolitan`, `govet`, `loggercheck`, `makezero`, `musttag`, `nilerr`, `nilnesserr`, `noctx`, `protogetter`, `reassign`, `recvcheck`, `rowserrcheck`, `spancheck`, `sqlclosecheck`, `staticcheck`, `testifylint`, `zerologlint` |
| comment | `dupword`, `godot`, `godox`, `misspell` | 
| complexity | `cyclop`, `funlen`, `gocognit`, `gocyclo`, `maintidx`, `nestif` |
| error | `err113`, `errcheck`, `errorlint`, `wrapcheck` |
| format | `gci`, `gofmt`, `gofumpt`, `goimports` |
| import | `depguard`, `gci`, `goimports`, `gomodguard` |
| metalinter | `gocritic`, `govet`, `revive`, `staticcheck` |
| module | `depguard`, `gomoddirectives`, `gomodguard` |
| performance | `bodyclose`, `fatcontext`, `noctx`, `perfsprint`, `prealloc` |
| sql | `rowserrcheck`, `sqlclosecheck` |
| style | `asciicheck`, `canonicalheader`, `containedctx`, `copyloopvar`, `decorder`, `depguard`, `dogsled`, `dupl`, `err113`, `errname`, `exhaustruct`, `exptostd`, `forbidigo`, `forcetypeassert`, `ginkgolinter`, `gochecknoglobals`, `gochecknoinits`, `goconst`, `gocritic`, `godot`, `godox`, `goheader`, `gomoddirectives`, `gomodguard`, `goprintffuncname`, `gosimple`, `grouper`, `iface`, `importas`, `inamedparam`, `interfacebloat`, `intrange`, `ireturn`, `lll`, `loggercheck`, `makezero`, `mirror`, `misspell`, `mnd`, `musttag`, `nakedret`, `nilnil`, `nlreturn`, `nolintlint`, `nonamedreturns`, `nosprintfhostport`, `paralleltest`, `predeclared`, `promlinter`, `revive`, `sloglint`, `stylecheck`, `tagalign`, `tagliatelle`, `testpackage`, `tparallel`, `unconvert`, `usestdlibvars`, `varnamelen`, `wastedassign`, `whitespace`, `wrapcheck`, `wsl` |
| test | `exhaustruct`, `paralleltest`, `testableexamples`, `testifylint`, `testpackage`, `thelper`, `tparallel`, `usetesting` |
| unused | `ineffassign`, `unparam`, `unused` |

{{< /tab >}}
{{< tab >}}
```yaml
# Removed
```
{{< /tab >}}
{{< /tabs >}}

#### `typecheck`

This `typecheck` is not a linter, so it cannot be enabled or disabled:

- [FAQ: Why do you have typecheck errors?](/docs/welcome/faq/#why-do-you-have-typecheck-errors)
- [FAQ: Why is it not possible to skip/ignore typecheck errors?](/docs/welcome/faq/#why-is-it-not-possible-to-skipignore-typecheck-errors)

#### Deprecated Linters

The following deprecated linters have been removed:

- `deadcode`
- `execinquery`
- `exhaustivestruct`
- `exportloopref`
- `golint`
- `ifshort`
- `interfacer`
- `maligned`
- `nosnakecase`
- `scopelint`
- `structcheck`
- `tenv`
- `varcheck`

#### Alternative Linter Names

The alternative linters has been removed.

| Alt Name v1 | Name v2       |
|-------------|---------------|
| `gas`       | `gosec`       |
| `goerr113`  | `err113`      |
| `gomnd`     | `mnd`         |
| `logrlint`  | `loggercheck` |
| `megacheck` | `staticcheck` |
| `vet`       | `govet`       |
| `vetshadow` | `govet`       |

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  enable:
    - gas
    - goerr113
    - gomnd
    - logrlint
    - megacheck
    - vet
    - vetshadow
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  enable:
    - gosec
    - err113
    - mnd
    - loggercheck
    - staticcheck
    - govet
```
{{< /tab >}}
{{< /tabs >}}

### `linters-settings`

The `linters-settings` section has been split into `linters.settings` and `formatters.settings`.

Settings for `gci`, `gofmt`, `gofumpt`, and `goimports` are moved to the `formatters.settings` section.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  govet:
    enable-all: true
  gofmt:
    simplify: false
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    govet:
      enable-all: true

formatters:
  settings:
    gofmt:
      simplify: false
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.asasalint.ignore-test`

This option has been removed.

To ignore test files, use `linters.exclusions.rules`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  asasalint:
    ignore-test: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    rules:
      - path: '(.+)_test\.go'
        linters:
          - asasalint
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.copyloopvar.ignore-alias`

This option has been deprecated since v1.58.0 and has been replaced with `linters.settings.copyloopvar.check-alias`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  copyloopvar:
    ignore-alias: false
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    copyloopvar:
      check-alias: true
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.cyclop.skip-tests`

This option has been removed.

To ignore test files, use `linters.exclusions.rules`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  cyclop:
    skip-test: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    rules:
      - path: '(.+)_test\.go'
        linters:
          - cyclop
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.errcheck.exclude`

This option has been deprecated since v1.42.0 and has been removed.

To exclude functions, use `linters.settings.errcheck.exclude-functions` instead.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  errcheck:
    exclude: ./errcheck_excludes.txt
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    errcheck:
      exclude-functions:
        - io.ReadFile
        - io.Copy(*bytes.Buffer)
        - io.Copy(os.Stdout)
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.errcheck.ignore`

This option has been deprecated since v1.13.0 and has been removed.

To exclude functions, use `linters.settings.errcheck.exclude-functions` instead.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  errcheck:
    ignore: 'io:.*'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    errcheck:
      exclude-functions:
        - 'io.ReadFile'
        - 'io.Copy(*bytes.Buffer)'
        - 'io.Copy(os.Stdout)'
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.exhaustive.check-generated`

This option has been removed.

To analyze generated files, use `linters.exclusions.generated`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  exhaustive:
    check-generated: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    generated: disable
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.forbidigo.forbid[].p`

This field has been replaced with `linters-settings.forbidigo.forbid[].pattern`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  forbidigo:
    forbid:
      - p: '^fmt\.Print.*$'
        msg: Do not commit print statements.
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    forbidigo:
      forbid:
        - pattern: '^fmt\.Print.*$'
          msg: Do not commit print statements.
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.forbidigo.forbid[]<pattern>`

The `pattern` has become mandatory for the `forbid` field.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  forbidigo:
    forbid:
      - '^print(ln)?$'
      - '^spew\.(ConfigState\.)?Dump$'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    forbidigo:
      forbid:
        - pattern: '^print(ln)?$'
        - pattern: '^spew\.(ConfigState\.)?Dump$'
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.gci.local-prefixes`

This option has been deprecated since v1.44.0 and has been removed.

Use `linters.settings.gci.sections` instead.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  gci:
    local-prefixes: 'github.com/example/pkg'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    gci:
      sections:
        - standard
        - default
        - prefix(github.com/example/pkg)
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.gci.skip-generated`

This option has been removed.

To analyze generated files, use `linters.exclusions.generated`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  settings:
    gci:
      skip-generated: false
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    generated: disable
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.goconst.ignore-tests`

This option has been removed.

To ignore test files, use `linters.exclusions.rules`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  goconst:
    ignore-tests: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    rules:
      - path: '(.+)_test\.go'
        linters:
          - goconst
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.gocritic.settings.ruleguard.rules`

The special variable `${configDir}` has been replaced with `${base-path}`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  gocritic:
    settings:
      ruleguard:
        rules: '${configDir}/ruleguard/rules-*.go'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    gocritic:
      settings:
        ruleguard:
          rules: '${base-path}/ruleguard/rules-*.go'
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.govet.check-shadowing`

This option has been deprecated since v1.57.0 and has been removed.

Use `linters.settings.govet.enable: shadow` instead.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  govet:
    check-shadowing: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    govet:
      enable:
        - shadow
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.misspell.ignore-words`

This option has been replaced with `linters.settings.misspell.ignore-rules`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  misspell:
    ignore-words:
      - foo
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    misspell:
      ignore-rules:
        - foo
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.predeclared.ignore`

This string option has been replaced with the slice option with the same name.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  predeclared:
    ignore: "new,int"
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    predeclared:
      ignore:
        - new
        - int
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.predeclared.q`

This option has been replaced with `linters.settings.predeclared.qualified-name`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  predeclared:
    q: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    predeclared:
      qualified-name: true
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.revive.ignore-generated-header`

This option has been removed.

Use `linters.exclusions.generated` instead.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  revive:
    ignore-generated-header: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    generated: strict
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.sloglint.context-only`

This option has been deprecated since v1.58.0 and has been replaced with `linters.settings.sloglint.context`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  sloglint:
    context-only: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    sloglint:
      context: all
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.staticcheck.go`

This option has been deprecated since v1.47.0 and has been removed.

Use `run.go` instead.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  staticcheck:
    go: '1.22'
```
{{< /tab >}}
{{< tab >}}
```yaml
run:
  go: '1.22'
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.unused.exported-is-used`

This option has been deprecated since v1.60.0 and has been removed.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  unused:
    exported-is-used: true
```
{{< /tab >}}
{{< tab >}}
```yaml
# Removed
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.usestdlibvars.os-dev-null`

This option has been deprecated since v1.51.0 and has been removed.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  usestdlibvars:
    os-dev-null: true
```
{{< /tab >}}
{{< tab >}}
```yaml
# Removed
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.usestdlibvars.syslog-priority`

This option has been deprecated since v1.51.0 and has been removed.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  usestdlibvars:
    syslog-priority: true
```
{{< /tab >}}
{{< tab >}}
```yaml
# Removed
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.wrapcheck.ignoreInterfaceRegexps`

This option has been renamed to `linters.settings.wrapcheck.ignore-interface-regexps`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  wrapcheck:
    ignoreInterfaceRegexps:
      - '^(?i)c(?-i)ach(ing|e)'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    wrapcheck:
      ignore-interface-regexps:
        - '^(?i)c(?-i)ach(ing|e)'
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.wrapcheck.ignorePackageGlobs`

This option has been renamed to `linters.settings.wrapcheck.ignore-package-globs`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  wrapcheck:
    ignorePackageGlobs:
      - 'encoding/*'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    wrapcheck:
      ignore-package-globs:
        - 'encoding/*'
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.wrapcheck.ignoreSigRegexps`

This option has been renamed to `linters.settings.wrapcheck.ignore-sig-regexps`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
  linters-settings:
    wrapcheck:
      ignoreSigRegexps:
        - '\.New.*Error\('
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    wrapcheck:
      ignore-sig-regexps:
        - '\.New.*Error\('
```
{{< /tab >}}
{{< /tabs >}}

#### `linters-settings.wrapcheck.ignoreSigs`

This option has been renamed to `linters.settings.wrapcheck.ignore-sigs`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters-settings:
  wrapcheck:
    ignoreSigs:
      - '.Errorf('
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  settings:
    wrapcheck:
      ignore-sigs:
        - '.Errorf('
```
{{< /tab >}}
{{< /tabs >}}

### `issues`

#### `issues.exclude-case-sensitive`

This property has been removed.

`issues.exclude`, `issues.exclude-rules.text`, and `issues.exclude-rules.source` are case-sensitive by default.

To ignore case, use `(?i)` at the beginning of a regex syntax.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude-case-sensitive: false
  exclude:
    - 'abcdef'
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    rules:
      - path: '(.+)\.go$'
        text: (?i)abcdef
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude-dirs-use-default`

This property has been removed.

Use `linters.exclusions.paths` and `formatters.exclusions.paths` to exclude directories.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude-dirs-use-default: true
```
{{< /tab >}}
{{< tab >}}

```yaml
linters:
  exclusions:
    paths:
      - third_party$
      - builtin$
      - examples$
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude-dirs`

This property has been replaced with `linters.exclusions.paths` and `formatters.exclusions.paths`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude-dirs:
    - src/external_libs
    - autogenerated_by_my_lib
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    paths:
      - src/external_libs
      - autogenerated_by_my_lib
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude-files`

This property has been replaced with `linters.exclusions.paths` and `formatters.exclusions.paths`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude-files:
    - '.*\.my\.go$'
    - lib/bad.go
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    paths:
      - '.*\.my\.go$'
      - lib/bad.go
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude-generated-strict`

This property has been deprecated since v1.59.0 and has been replaced with `linters.exclusions.generated: strict`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  exclude-generated-strict: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    generated: strict
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude-generated`

This property has been replaced with `linters.exclusions.generated`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
linters:
  exclude-generated: lax
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    generated: lax
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude-rules`

This property has been replaced with `linters.exclusions.rules`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude-rules:
    - path: '_test\.go'
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
    - path-except: '_test\.go'
      linters:
        - staticcheck
    - path: internal/hmac/
      text: "weak cryptographic primitive"
      linters:
        - gosec
    - linters:
        - staticcheck
      text: "SA9003:"
    - linters:
        - err113
      source: "foo"
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    rules:
      - path: '_test\.go'
        linters:
          - dupl
          - errcheck
          - gocyclo
          - gosec
      - path-except: '_test\.go'
        linters:
          - staticcheck
      - path: internal/hmac/
        text: weak cryptographic primitive
        linters:
          - gosec
      - text: 'SA9003:'
        linters:
          - staticcheck
      - source: foo
        linters:
          - err113
  ```
{{< /tab >}}
{{< /tabs >}}


#### `issues.exclude-use-default`

This property has been replaced with `linters.exclusions.presets`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude-use-default: true
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    presets:
      - comments
      - common-false-positives
      - legacy
      - std-error-handling
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.exclude`

This property has been replaced with `linters.exclusions.rules`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  exclude:
    - abcdef
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
  rules:
    - path: '(.+)\.go$'
      text: abcdef
```
{{< /tab >}}
{{< /tabs >}}

#### `issues.include`

This property has been replaced with `linters.exclusions.presets`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
issues:
  include:
    - EXC0014
    - EXC0015
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    presets:
      - common-false-positives
      - legacy
      - std-error-handling
```
{{< /tab >}}
{{< /tabs >}}

### `output`

#### `output.format`

This property has been deprecated since v1.57.0 and has been replaced with `output.formats`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  format: 'checkstyle:report.xml,json:stdout,colored-line-number'
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    checkstyle:
      path: 'report.xml'
    json:
      path: stdout
    text:
      path: stdout
      color: true
```
{{< /tab >}}
{{< /tabs >}}

#### `output.formats[].format: <name>`

The property `output.formats[].format` has been replaced with `output.formats[].<format_name>`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  formats:
    - format: json
      path: stderr
    - format: checkstyle
      path: report.xml
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    json:
      path: stderr
    checkstyle:
      path: report.xml
```
{{< /tab >}}
{{< /tabs >}}

#### `output.formats[].format: line-number`

This format has been replaced by the format `text`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  formats:
    - format: line-number
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    text:
      path: stdout
```
{{< /tab >}}
{{< /tabs >}}

#### `output.formats[].format: colored-line-number`

This format has been replaced by the format `text` with the option `colors` (`true` by default).

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  formats:
    - format: colored-line-number
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    text:
      path: stdout
      colors: true
```
{{< /tab >}}
{{< /tabs >}}

#### `output.formats[].format: colored-tab`

This format has been replaced by the format `tab` with the option `colors` (`true` by default).

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  formats:
    - format: colored-tab
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    tab:
      path: stdout
      colors: true
```
{{< /tab >}}
{{< /tabs >}}

#### `output.print-issued-lines`

This property has been removed.

To not print the lines with issues, use the `text` format with the option `print-issued-lines: false`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  formats:
    - format: line-number
      path: stdout
  print-issued-lines: false
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    text:
      path: stdout
      print-issued-lines: false
```
{{< /tab >}}
{{< /tabs >}}

#### `output.print-linter-name`

This property has been removed.

To not print the linter name, use the `text` format with the option `print-linter-name: false`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  formats:
    - format: line-number
      path: stdout
  print-linter-name: false
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  formats:
    text:
      path: stdout
      print-linter-name: false
```
{{< /tab >}}
{{< /tabs >}}

#### `output.show-stats`

This property is `true` by default.

#### `output.sort-order`

This property has a new default value `['linter', 'file']` instead of `['file']`.

#### `output.sort-results`

The property has been removed.

The output results are always sorted.

#### `output.uniq-by-line`

This property has been deprecated since v1.63.0 and has been replaced by `issues.uniq-by-line`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
output:
  uniq-by-line: true
```
{{< /tab >}}
{{< tab >}}
```yaml
issues:
  uniq-by-line: true
```
{{< /tab >}}
{{< /tabs >}}

### `run`

#### `run.go`

The new fallback value for this property is `1.22` instead of `1.17`.

#### `run.concurrency`

This property value set to match Linux container CPU quota by default and fallback on the number of logical CPUs in the machine.

#### `run.relative-path-mode`

This property has a new default value of `cfg` instead of `wd`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
run:
# When not specified, relative-path-mode is set to 'wd' by default
```
{{< /tab >}}
{{< tab >}}
```yaml
run:
  relative-path-mode: 'cfg'
```
{{< /tab >}}
{{< /tabs >}}

#### `run.show-stats`

This property has been deprecated since v1.57.0 and has been replaced by `output.show-stats`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
run:
  show-stats: true
```
{{< /tab >}}
{{< tab >}}
```yaml
output:
  show-stats: true
```
{{< /tab >}}
{{< /tabs >}}

#### `run.skip-dirs-use-default`

This property has been deprecated since v1.57.0 and has been replaced by `issues.exclude-dirs-use-default`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
run:
  skip-dirs-use-default: false
```
{{< /tab >}}
{{< tab >}}
```yaml
issues:
  exclude-dirs-use-default: false
```
{{< /tab >}}
{{< /tabs >}}

#### `run.skip-dirs`

This property has been deprecated since v1.57.0 and has been removed.

Use `linters.exclusions.paths` and `formatters.exclusions.paths` to exclude directories.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
run:
  skip-dirs:
    - src/external_libs
    - autogenerated_by_my_lib
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    paths:
      - src/external_libs
      - autogenerated_by_my_lib
```
{{< /tab >}}
{{< /tabs >}}

#### `run.skip-files`

This property has been deprecated since v1.57.0 and has been removed.

Use `linters.exclusions.paths` and `formatters.exclusions.paths` to exclude files.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
run:
  skip-files:
    - '.*\.my\.go$'
    - lib/bad.go
```
{{< /tab >}}
{{< tab >}}
```yaml
linters:
  exclusions:
    paths:
      - '.*\.my\.go$'
      - lib/bad.go
```
{{< /tab >}}
{{< /tabs >}}

#### `run.timeout`

This property value is disabled by default (`0`).

### `severity`

#### `severity.default-severity`

This property has been replaced with `severity.default`.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
severity:
  default-severity: error
```
{{< /tab >}}
{{< tab >}}
```yaml
severity:
  default: error
```
{{< /tab >}}
{{< /tabs >}}

#### `severity.rules.case-sensitive`

`severity.rules.text` and `severity.rules.source` are case-sensitive by default.

To ignore case, use `(?i)` at the beginning of a regex syntax.

{{< tabs items="v1,v2" >}}
{{< tab >}}
```yaml
severity:
  case-sensitive: true
  rules:
    - severity: info
      linters:
        - foo
      text: 'Example.*'
```
{{< /tab >}}
{{< tab >}}
```yaml
severity:
  rules:
    - severity: info
      linters:
        - foo
      text: '(?i)Example.*'
```
{{< /tab >}}
{{< /tabs >}}

### `version`

The `version` property has been added to the configuration file.

```yaml
version: "2"
```

### Integration

#### VSCode

{{< tabs items="v1,v2" >}}
{{< tab >}}
```JSONata
"go.lintTool": "golangci-lint",
"go.lintFlags": [
"--fast"
]
```
{{< /tab >}}
{{< tab >}}
```JSONata
"go.lintTool": "golangci-lint",
"go.lintFlags": [
"--fast-only"
],
"go.formatTool": "custom",
"go.alternateTools": {
"customFormatter": "golangci-lint"
},
"go.formatFlags": [
"fmt",
"--stdin"
]
```
{{< /tab >}}
{{< /tabs >}}

### Command Line Flags

The following flags have been removed:

- `--disable-all`
- `--enable-all`
- `-p, --presets`
- `--fast`
- `-e, --exclude`
- `--exclude-case-sensitive`
- `--exclude-dirs-use-default`
- `--exclude-dirs`
- `--exclude-files`
- `--exclude-generated`
- `--exclude-use-default`
- `--go string`
- `--sort-order`
- `--sort-results`
- `--out-format`
- `--print-issued-lines`
- `--print-linter-name`

#### `--out-format`

`--out-format` has been replaced with the following flags:

```bash
# Previously 'colored-line-number' and 'line-number'
--output.text.path
--output.text.print-linter-name
--output.text.print-issued-lines
--output.text.colors
```

```bash
# Previously 'json'
--output.json.path
```

```bash
# Previously 'colored-tab' and 'tab'
--output.tab.path
--output.tab.print-linter-name
--output.tab.colors
```

```bash
# Previously 'html'
--output.html.path
```

```bash
# Previously 'checkstyle'
--output.checkstyle.path
```

```bash
# Previously 'code-climate'
--output.code-climate.path
```

```bash
# Previously 'junit-xml' and 'junit-xml-extended'
--output.junit-xml.path
--output.junit-xml.extended
```

```bash
# Previously 'teamcity'
--output.teamcity.path
```

```bash
# Previously 'sarif'
--output.sarif.path
```

#### `--print-issued-lines`

`--print-issued-lines` has been replaced with `--output.text.print-issued-lines`.

#### `--print-linter-name`

`--print-linter-name` has been replaced with `--output.text.print-linter-name` or `--output.tab.print-linter-name`.

#### `--disable-all` and `--enable-all`

`--disable-all` has been replaced with `--default=none`.

`--enable-all` has been replaced with `--default=all`.

#### Examples

Run only the `govet` linter, output results to stdout in JSON format, and sort results:

{{< tabs items="v1,v2" >}}
{{< tab >}}

```bash
golangci-lint run --disable-all --enable=govet --out-format=json --sort-order=linter --sort-results
```
{{< /tab >}}
{{< tab >}}

```bash
golangci-lint run --default=none --enable=govet --output.json.path=stdout
```
{{< /tab >}}
{{< /tabs >}}

Do not print issued lines, output results to stdout without colors in text format, and to `gl-code-quality-report.json` file in Code Climate's format:

{{< tabs items="v1,v2" >}}
{{< tab >}}

```bash
golangci-lint run --print-issued-lines=false --out-format code-climate:gl-code-quality-report.json,line-number
```
{{< /tab >}}
{{< tab >}}

```bash
golangci-lint run --output.text.path=stdout --output.text.colors=false --output.text.print-issued-lines=false --output.code-climate.path=gl-code-quality-report.json
```
{{< /tab >}}
{{< /tabs >}}
