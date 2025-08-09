---
title: Command-Line
weight: 1
---

```console
$ golangci-lint -h
{.CmdRootHelpText}
```

## `run`

{{< cards >}}
    {{< card link="/docs/linters" title="Linters Overview" icon="collection" >}}
    {{< card link="/docs/configuration/file/#linters-configuration" title="Global Configuration" icon="adjustments" >}}
    {{< card link="/docs/linters/configuration/" title="Linter Settings" icon="adjustments" >}}
{{< /cards >}}

```console
$ golangci-lint run -h
{.CmdRunHelpText}
```

When the `--cpu-profile-path` or `--mem-profile-path` arguments are specified,
golangci-lint writes runtime profiling data in the format expected by the [pprof](https://github.com/google/pprof) visualization tool.

When the `--trace-path` argument is specified, `golangci-lint` writes runtime tracing data in the format expected by
the `go tool trace` command and visualization tool.

## fmt

{{< cards >}}
{{< card link="/docs/formatters" title="Formatters Overview" icon="collection" >}}
{{< card link="/docs/configuration/file/#formatters-configuration" title="Global Configuration" icon="adjustments" >}}
{{< card link="/docs/formatters/configuration/" title="Formatter Settings" icon="adjustments" >}}
{{< /cards >}}

```console
$ golangci-lint fmt -h
{.CmdFmtHelpText}
```

## `migrate`

```console
$ golangci-lint migrate -h
{.CmdMigrateHelpText}
```

## `formatters`

```console
$ golangci-lint formatters -h
{.CmdFormattersHelpText}
```

## `help`

```console
$ golangci-lint help -h
{.CmdHelpText}
```

## `linters`

```console
$ golangci-lint linters -h
{.CmdLintersHelpText}
```

## `cache`

Golangci-lint stores its cache in the subdirectory `golangci-lint` inside the [default user cache directory](https://pkg.go.dev/os#UserCacheDir).

You can override the default cache directory with the environment variable `GOLANGCI_LINT_CACHE`; the path must be absolute.

The cache is only used by `golangci-lint run` (linters).

```console
$ golangci-lint cache -h
{.CmdCacheHelpText}
```

## `config`

```console
$ golangci-lint config -h
{.CmdConfigHelpText}
```

## `custom`

```console
$ golangci-lint custom -h
{.CmdCustomHelpText}
```

## `version`

```console
$ golangci-lint version -h
{.CmdVersionHelpText}
```

## `completion`

```console
$ golangci-lint completion -h
{.CmdCompletionHelpText}
```
