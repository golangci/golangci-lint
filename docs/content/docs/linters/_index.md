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

## Enabled by Default

{.EnabledByDefaultLinters}

## Disabled by Default

{.DisabledByDefaultLinters}
