---
title: Formatters
weight: 4
excludeSearch: true
aliases:
  - /usage/formatters/
---

To see a list of supported formatters and which formatters are enabled/disabled:

```bash
golangci-lint help formatters
```

To see a list of formatters enabled by your configuration, use:

```bash
golangci-lint formatters
```

{{< cards cols=2 >}}
    {{< card link="/docs/welcome/quick-start/#formatting" title="Quick Start" icon="terminal" >}}
    {{< card link="/docs/configuration/cli/#fmt" title="CLI" icon="terminal" >}}
    {{< card link="/docs/configuration/file/#formatters-configuration" title="Global Configuration" icon="adjustments" >}}
    {{< card link="/docs/formatters/configuration/" title="Formatter Settings" icon="adjustments" >}}
{{< /cards >}}

## All formatters

{{< golangci/items/filter >}}
    {{< golangci/items/filter-badge class="gl-filter" data="new" icon="fire" content="New" color="yellow" >}}
    {{< golangci/items/filter-badge class="gl-filter" data="deprecated" icon="emoji-sad" content="Deprecated" color="blue" >}}
    {{< golangci/items/filter-badge class="gl-filter-reset gl-hidden" icon="trash" content="Reset" color="red" border=true >}}
{{< /golangci/items/filter >}}

{{< cards >}}
    {{< golangci/items/cards path="formatters" data="formatters_info" >}}
{{< /cards >}}
