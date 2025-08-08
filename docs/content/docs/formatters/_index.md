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

<div class="hx-mt-6">
    {{< icon "filter" >}}
    {{< clickable-badge icon="fire" id="new-filter" content="New" class="gl-filter hx-cursor-pointer" type="warning" border=false >}}
    {{< clickable-badge icon="emoji-sad" id="deprecated-filter" content="Deprecated" class="gl-filter hx-cursor-pointer" type="info" border=false >}}
    {{< clickable-badge icon="trash" id="reset-filter" content="Reset" class="gl-filter-reset gl-hidden hx-cursor-pointer" type="error" border=true >}}
</div>

{{< cards >}}
    {{< item-cards path="formatters" data="formatters_info" >}}
{{< /cards >}}
