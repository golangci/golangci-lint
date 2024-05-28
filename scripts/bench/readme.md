# Benchmarks

The script use [Hyperfine](https://github.com/sharkdp/hyperfine) to benchmark the command line of golangci-lint.

## Benchmark one linter: with a local version

```bash
make bench_local LINTER=gosec VERSION=v1.59.0
```

## Benchmark one linter: between 2 existing versions

```bash
make bench_version LINTER=gosec VERSION_OLD=v1.58.1 VERSION_NEW=v1.59.0 
```
