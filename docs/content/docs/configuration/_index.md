---
title: Configuration
weight: 2
---

The config file has lower priority than command-line options.
If the same bool/string/int option is provided on the command-line
and in the config file, the option from command-line will be used.
Slice options (e.g. list of enabled/disabled linters) are combined from the command-line and config file.

## More

{{< cards cols=2 >}}
  {{< card link="/docs/configuration/cli" title="Command Line" icon="terminal" >}}
  {{< card link="/docs/configuration/file" title="Configuration File" icon="adjustments" >}}
{{< /cards >}}


## Cache

Golangci-lint stores its cache in the subdirectory `golangci-lint` inside the [default user cache directory](https://pkg.go.dev/os#UserCacheDir).

You can override the default cache directory with the environment variable `GOLANGCI_LINT_CACHE`; the path must be absolute.
