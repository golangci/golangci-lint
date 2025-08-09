---
title: 'Golangci-lint'
layout: hextra-home
params:
  width: wide
---

{{< hextra/hero-container image="images/golangci-lint-logo-anim.gif" imageWidth="300" imageHeight="300" imageTitle="golangci-lint" >}}

<div class="hx-mt-6 hx-mb-6">
{{< hextra/hero-headline >}}
  Golangci-lint is a fast linters runner for Go
{{< /hextra/hero-headline >}}
</div>

<div class="hx-mb-12">
{{< hextra/hero-subtitle >}}
  It runs linters in parallel, uses caching, supports YAML configuration,&nbsp;<br class="sm:hx-block hx-hidden" />integrates with all major IDEs, and includes over a hundred linters.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx-mb-6">
{{< hextra/hero-button text="Get Started" link="docs" >}}
</div>

{{< /hextra/hero-container >}}

<div class="hx-mt-6"></div>

{{< hextra/feature-grid cols=3 >}}
  {{< hextra/feature-card
    icon="fast-forward"
    title="Fast"
    subtitle="Runs linters in parallel, reuses Go build cache and caches analysis results."
    style="background: radial-gradient(ellipse at 50% 80%,rgba(194,97,254,0.15),hsla(0,0%,100%,0));"
    link="/docs/" >}}
  {{< hextra/feature-card
    icon="desktop-computer"
    title="Integrations"
    subtitle="Integrations with VS Code, Sublime Text, GoLand, GNU Emacs, Vim, GitHub Actions."
    style="background: radial-gradient(ellipse at 50% 80%,rgba(142,53,74,0.15),hsla(0,0%,100%,0));"
    link="/docs/welcome/integrations" >}}
  {{< hextra/feature-card
    icon="sparkles"
    title="Nice outputs"
    subtitle="Text with colors and source code lines, JSON, tab, HTML, Checkstyle, Code-Climate, JUnit-XML, TeamCity, SARIF."
    style="background: radial-gradient(ellipse at 50% 80%,rgba(221,210,59,0.15),hsla(0,0%,100%,0));"
    link="/docs/configuration/file/#output-configuration" >}}
  {{< hextra/feature-card
    icon="eye-off"
    title="Minimum number of false positives"
    subtitle="Tuned default settings."
    link="/docs/linters/false-positives" >}}
  {{< hextra/feature-card
    icon="collection"
    title="A lot of linters"
    subtitle="No need to install them."
    link="/docs/linters" >}}
  {{< hextra/feature-card
    icon="document-text"
    title="YAML-based configuration"
    subtitle="Easy to read and maintain."
    link="/docs/configuration/file" >}}
{{< /hextra/feature-grid >}}
