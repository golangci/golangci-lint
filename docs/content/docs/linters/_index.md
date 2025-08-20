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
    {{< golangci/items/clickable-badge class="gl-filter" id="default-filter" icon="inbox" content="Default" type="info" >}}
    {{< golangci/items/clickable-badge class="gl-filter" id="new-filter" icon="fire" content="New" type="warning" >}}
    {{< golangci/items/clickable-badge class="gl-filter" id="autofix-filter" icon="sparkles" content="Autofix" type="info" >}}
    {{< golangci/items/clickable-badge class="gl-filter" id="fast-filter" icon="fast-forward" content="Fast" >}}
    {{< golangci/items/clickable-badge class="gl-filter" id="slow-filter" icon="play" content="Slow" >}}
    {{< golangci/items/clickable-badge class="gl-filter" id="deprecated-filter" icon="emoji-sad" content="Deprecated" type="info" >}}
    {{< golangci/items/clickable-badge class="gl-filter-reset gl-hidden" icon="trash" id="reset-filter" content="Reset" type="error" border=true >}}
{{< /golangci/items/filter >}}

{{< cards >}}
    {{< golangci/items/cards path="linters" data="linters_info" >}}
{{< /cards >}}
