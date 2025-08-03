---
title: Command-Line
weight: 1
---

## run

```sh
$ golangci-lint run -h
{.RunHelpText}
```

When the `--cpu-profile-path` or `--mem-profile-path` arguments are specified,
golangci-lint writes runtime profiling data in the format expected by the [pprof](https://github.com/google/pprof) visualization tool.

When the `--trace-path` argument is specified, `golangci-lint` writes runtime tracing data in the format expected by
the `go tool trace` command and visualization tool.

## fmt

```sh
$ golangci-lint fmt -h
{.FmtHelpText}
```
