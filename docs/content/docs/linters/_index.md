---
title: Linters
weight: 3
excludeSearch: true
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

<div class="hx-mt-6">
    {{< icon "filter" >}}
    {{< clickable-badge icon="inbox" id="default-filter" content="Default" class="gl-filter hx-cursor-pointer" type="info" border=false >}}
    {{< clickable-badge icon="fire" id="new-filter" content="New" class="gl-filter hx-cursor-pointer" type="warning" border=false >}}
    {{< clickable-badge icon="sparkles" id="autofix-filter" content="Autofix" class="gl-filter hx-cursor-pointer" type="info" border=false >}}
    {{< clickable-badge icon="fast-forward" id="fast-filter" content="Fast" class="gl-filter hx-cursor-pointer" border=false >}}
    {{< clickable-badge icon="play" id="slow-filter" content="Slow" class="gl-filter hx-cursor-pointer" border=false >}}
    {{< clickable-badge icon="emoji-sad" id="deprecated-filter" content="Deprecated" class="gl-filter hx-cursor-pointer" type="info" border=false >}}
    {{< clickable-badge icon="trash" id="reset-filter" content="Reset" class="gl-filter-reset gl-hidden hx-cursor-pointer" type="error" border=true >}}
</div>

{{< cards >}}
    {{< item-cards path="linters" data="linters_info" >}}
{{< /cards >}}
