---
title: Configuration File
weight: 2
---

Golangci-lint looks for config files in the following paths from the current working directory:

- `.golangci.yml`
- `.golangci.yaml`
- `.golangci.toml`
- `.golangci.json`

Golangci-lint also searches for config files in all directories from the directory of the first analyzed path up to the root.
If no configuration file has been found, golangci-lint will try to find one in your home directory.
To see which config file is being used and where it was sourced from run golangci-lint with `-v` option.

Config options inside the file are identical to command-line options.
You can configure specific linters' options only within the config file (not the command-line).

There is a [`.golangci.reference.yml`](https://github.com/golangci/golangci-lint/blob/HEAD/.golangci.reference.yml) file with all supported options, their descriptions, and default values.
This file is neither a working example nor a recommended configuration,
it's just a reference to display all the configuration options used to generate the documentation.

The configuration file can be validated with the JSON Schema: [golangci.jsonschema.json](https://golangci-lint.run/jsonschema/golangci.jsonschema.json)

{{% configuration-file-snippet section="root" %}}

## `version` configuration

{{% configuration-file-snippet section="version" %}}

## `linters` configuration

{{< cards  cols=2 >}}
{{< card link="/docs/linters" title="Linters Overview" icon="collection" >}}
{{< card link="/docs/linters/configuration" title="Linters  Settings" icon="adjustments" >}}
{{< /cards >}}

{{% configuration-file-snippet section="linters" %}}

## `formatters` configuration

{{< cards  cols=2 >}}
{{< card link="/docs/formatters" title="Formatters Overview" icon="collection" >}}
{{< card link="/docs/formatters/configuration" title="Formatters  Settings" icon="adjustments" >}}
{{< /cards >}}

{{% configuration-file-snippet section="formatters" %}}

## `issues` configuration

{{% configuration-file-snippet section="issues" %}}

## `output` configuration

{{% configuration-file-snippet section="output" %}}

## `run` configuration

{{% configuration-file-snippet section="run" %}}

## `severity` configuration

{{% configuration-file-snippet section="severity" %}}
