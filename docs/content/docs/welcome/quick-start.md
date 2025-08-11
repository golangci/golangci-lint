---
title: Quick Start
weight: 2
aliases:
  - /welcome/quick-start/
---

## Linting

To run golangci-lint:

```bash
golangci-lint run
```

It's an equivalent of:

```bash
golangci-lint run ./...
```

You can choose which directories or files to analyze:

```bash
golangci-lint run dir1 dir2/...
golangci-lint run file1.go
```

Directories are NOT analyzed recursively.
To analyze them recursively append `/...` to their path.
It's not possible to mix files and packages/directories, and files must come from the same package.

Golangci-lint can be used with zero configuration. By default, the following linters are enabled:

{{% cli-output section="defaultEnabledLinters" cmd="help linters" %}}

Pass `-E/--enable` to enable linter and `-D/--disable` to disable:

```bash
golangci-lint run --default=none -E errcheck
```

More information about available linters can be found in the [linters page](/docs/linters/).

## Formatting

To format your code:

```bash
golangci-lint fmt
```

You can choose which directories or files to analyze:

```bash
golangci-lint fmt dir1 dir2/...
golangci-lint fmt file1.go
```

More information about available formatters can be found in the [formatters page](/docs/formatters/).
