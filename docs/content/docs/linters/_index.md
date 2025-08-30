---
title: Linters
weight: 3
excludeSearch: true
aliases:
  - /usage/linters/
---

To see a list of supported linters and which linters are enabled/disabled:

```bash
golangci-lint help linters
```

To see a list of linters enabled by your configuration, use:

```bash
golangci-lint linters
```

{{< cards cols=2 >}}
    {{< card link="/docs/welcome/quick-start/#linting" title="Quick Start" icon="terminal" >}}
    {{< card link="/docs/configuration/cli/#run" title="CLI" icon="terminal" >}}
    {{< card link="/docs/configuration/file/#linters-configuration" title="Global Configuration" icon="adjustments" >}}
    {{< card link="/docs/linters/configuration/" title="Linter Settings" icon="adjustments" >}}
{{< /cards >}}

## All Linters

{{< golangci/items/filter >}}
    {{< golangci/items/filter-badge class="gl-filter" data="default" icon="inbox" content="Default" color="indigo" >}}
    {{< golangci/items/filter-badge class="gl-filter" data="new" icon="fire" content="New" color="yellow" >}}
    {{< golangci/items/filter-badge class="gl-filter" data="autofix" icon="sparkles" content="Autofix" color="blue" >}}
    {{< golangci/items/filter-badge class="gl-filter" data="fast" icon="fast-forward" content="Fast" >}}
    {{< golangci/items/filter-badge class="gl-filter" data="slow" icon="play" content="Slow" >}}
    {{< golangci/items/filter-badge class="gl-filter" data="deprecated" icon="emoji-sad" content="Deprecated" color="blue" >}}
    {{< golangci/items/filter-badge class="gl-filter-reset gl-hidden" icon="trash" content="Reset" color="red" border=true >}}
{{< /golangci/items/filter >}}

{{< cards >}}
    {{< golangci/items/cards path="linters" data="linters_info" >}}
{{< /cards >}}
