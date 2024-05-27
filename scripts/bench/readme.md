# Benchmarks

The script use [Hyperfine](https://github.com/sharkdp/hyperfine) to benchmark the command line of golangci-lint.

## Benchmark one linter: with a local version

```bash
LINTER=gosec VERSION=v1.59.0 make bench_local
```

## Benchmark one linter: between 2 existing versions

```bash
LINTER=gosec VERSION_OLD=v1.58.2 VERSION_NEW=v1.59.0 make bench_version
```
