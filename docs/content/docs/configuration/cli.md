---
title: Command-Line
weight: 1
---

{{% cli-output %}}

## `run`

{{< cards >}}
    {{< card link="/docs/linters" title="Linters Overview" icon="collection" >}}
    {{< card link="/docs/configuration/file/#linters-configuration" title="Global Configuration" icon="adjustments" >}}
    {{< card link="/docs/linters/configuration/" title="Linter Settings" icon="adjustments" >}}
{{< /cards >}}

> [!NOTE]
> This command executes enabled linters, and the formatters defined in [`formatters`](/docs/configuration/file/#formatters-configuration),
> but it does not format the code.
> 
> To only format code, use [`golangci-lint fmt`](/docs/configuration/cli/#fmt).
> To apply both linter fixes and formatting, use `golangci-lint run --fix`. 
> 
> The formatters cannot be enabled or disabled inside the [`linters`](/docs/configuration/file/#linters-configuration) section or the flags `-E/--enable`, `-D/--disable` of the command  [`golangci-lint run`](/docs/configuration/cli/#run).
> 
> The formatters can be enabled/disabled by defining them inside the [`formatters`](/docs/configuration/file/#formatters-configuration) section or by using the flags `-E/--enable`, `-D/--disable` of command [`golangci-lint fmt`](/docs/configuration/cli/#fmt).

{{% cli-output cmd="run" %}}

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

{{% cli-output cmd="fmt" %}}

## `migrate`

{{% cli-output cmd="migrate" %}}

## `formatters`

{{% cli-output cmd="formatters" %}}

## `help`

{{% cli-output cmd="help" %}}

## `linters`

{{% cli-output cmd="linters" %}}

## `cache`

Golangci-lint stores its cache in the subdirectory `golangci-lint` inside the [default user cache directory](https://pkg.go.dev/os#UserCacheDir).

You can override the default cache directory with the environment variable `GOLANGCI_LINT_CACHE`; the path must be absolute.

The cache is only used by `golangci-lint run` (linters).

{{% cli-output cmd="cache" %}}

## `config`

{{% cli-output cmd="config" %}}

## `custom`

{{% cli-output cmd="custom" %}}

## `version`

{{% cli-output cmd="version" %}}

## `completion`

{{% cli-output cmd="completion" %}}
