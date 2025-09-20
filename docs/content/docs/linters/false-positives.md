---
title: False Positives
weight: 3
aliases:
  - /usage/false-positives/
---

False positives are inevitable, but we did our best to reduce their count.

If a false positive occurred, you have several choices.

## Specific Linter Excludes

Most of the linters have a configuration, sometimes false-positives can be related to a bad configuration of a linter.
So it's recommended to check the linter configurations.

Otherwise, some linters have dedicated configuration to exclude or disable rules.

An example with `staticcheck`:

```yaml
linters:
  settings:
    staticcheck:
      checks:
        - all
        - '-SA1000' # disable the rule SA1000
        - '-SA1004' # disable the rule SA1004
```

## Exclude

### Exclude Issue by Text

You can use `linters.exclusions.rules` config option for per-path or per-linter configuration.

In the following example, all the reports from the linters (`linters`) that contains the text (`text`) are excluded:

```yaml
linters:
  exclusions:
    rules:
      - linters:
          - mnd
        text: "Magic number: 9"
```

In the following example, all the reports from the linters (`linters`) that originated from the source (`source`) are excluded:

```yaml
linters:
  exclusions:
    rules:
      - linters:
          - lll
        source: "^//go:generate "
```

In the following example, all the reports that contains the text (`text`) in the path (`path`) are excluded:

```yaml
linters:
  exclusions:
    rules:
      - path: path/to/a/file.go
        text: "string `example` has (\\d+) occurrences, make it a constant"
```

### Exclude Issues by Path

Exclude issues in path by `linters.exclusions.paths` or `linters.exclusions.rules` config options.

In the following example, all the reports from the linters (`linters`) that concerns the path (`path`) are excluded:

```yaml
linters:
  exclusions:
    rules:
      - path: '(.+)_test\.go'
        linters:
          - funlen
          - goconst
```

The opposite, excluding reports **except** for specific paths, is also possible.
In the following example, only test files get checked:

```yaml
linters:
  exclusions:
    rules:
      - path-except: '(.+)_test\.go'
        linters:
          - funlen
          - goconst
```

In the following example, all the reports related to the files (`paths`) are excluded:

```yaml
linters:
  exclusions:
    paths:
      - path/to/a/file.go
```

In the following example, all the reports related to the directories (`paths`) are excluded:

```yaml
linters:
  exclusions:
    paths:
      - path/to/a/dir/
```

## Nolint Directive

To exclude issues from all linters use `//nolint:all`.
For example, if it's used inline (not from the beginning of the line) it excludes issues only for this line.

```go
var bad_name int //nolint:all
```

To exclude issues from specific linters only:

```go
var bad_name int //nolint:wsl,unused
```

To exclude issues for the block of code, use this directive at the beginning of a line:

```go
//nolint:all
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

You may add a comment explaining or justifying why a `nolint` directive is being used on the same line as the flag itself:

```go
//nolint:gocyclo // This legacy function is complex, but the team too busy to simplify it
func someLegacyFunction() *string {
  // ...
}
```

You can see more examples of using `nolint` directives in [our tests](https://github.com/golangci/golangci-lint/tree/HEAD/pkg/result/processors/testdata) for it.

### Syntax

`nolint` is not regular comment but a directive.

> A directive comment is a line matching the regular expression `//(line |extern |export |[a-z0-9]+:[a-z0-9])`.
> https://go.dev/doc/comment#syntax

This means that no spaces are allowed between:
- `//` and `nolint`
- `nolint` and `:`
- `:` and the name of the linter.

Invalid syntax:
```go
// nolint
// nolint:xxx
//nolint :xxx
//nolint: xxx
```

Valid syntax:
```go
//nolint:xxx
```

## Exclusion Presets

Some exclusions are considered common.
To help golangci-lint users, those common exclusions are provided through presets.

{{% golangci/exclusion-presets-snippet %}}

{{% golangci/exclusion-preset-tables %}}
