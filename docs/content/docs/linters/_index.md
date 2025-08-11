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

{{< filter-bar >}}
    {{< clickable-badge class="gl-filter" id="default-filter" icon="inbox" content="Default" type="info" >}}
    {{< clickable-badge class="gl-filter" id="new-filter" icon="fire" content="New" type="warning" >}}
    {{< clickable-badge class="gl-filter" id="autofix-filter" icon="sparkles" content="Autofix" type="info" >}}
    {{< clickable-badge class="gl-filter" id="fast-filter" icon="fast-forward" content="Fast" >}}
    {{< clickable-badge class="gl-filter" id="slow-filter" icon="play" content="Slow" >}}
    {{< clickable-badge class="gl-filter" id="deprecated-filter" icon="emoji-sad" content="Deprecated" type="info" >}}
    {{< clickable-badge class="gl-filter-reset gl-hidden" icon="trash" id="reset-filter" content="Reset" type="error" border=true >}}
{{< /filter-bar >}}

{{< cards >}}
    {{< item-cards path="linters" data="linters_info" >}}
{{< /cards >}}
