---
title: Formatters
weight: 4
excludeSearch: true
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

{{< filter-bar >}}
    {{< clickable-badge class="gl-filter" id="new-filter" icon="fire" content="New" type="warning" >}}
    {{< clickable-badge class="gl-filter" id="deprecated-filter" icon="emoji-sad" content="Deprecated" type="info" >}}
    {{< clickable-badge class="gl-filter-reset gl-hidden" id="reset-filter" icon="trash" content="Reset" type="error" border=true >}}
{{< /filter-bar >}}

{{< cards >}}
    {{< item-cards path="formatters" data="formatters_info" >}}
{{< /cards >}}
