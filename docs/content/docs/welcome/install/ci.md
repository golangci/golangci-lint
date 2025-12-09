---
title: "CI Installation"
weight: 3
---

It's important to have reproducible CI: don't start to fail all builds at the same time.
With golangci-lint this can happen if you use option `linters.default: all` and a new linter is added
or even without `linters.default: all` when one upstream linter is upgraded.

> [!IMPORTANT]
> It's highly recommended installing a specific version of golangci-lint available on the [releases page](https://github.com/golangci/golangci-lint/releases).

## GitHub Actions

We recommend using [our GitHub Action](https://github.com/golangci/golangci-lint-action) for running golangci-lint in CI for GitHub projects.

It's [fast and uses smart caching](https://github.com/golangci/golangci-lint-action#performance) inside,
and it can be much faster than the simple binary installation.

Also, the action creates GitHub annotations for found issues (you don't need to dig into build log to see found by golangci-lint issues).

{{< cards cols=2 >}}
  {{< golangci/image-card src="/images/colored-line-number.png" title="Console Output" >}}
  {{< golangci/image-card src="/images/annotations.png" title="Annotations" >}}
{{< /cards >}}

## GitLab CI

GitLab provides a [guide for integrating golangci-lint into the Code Quality widget](https://docs.gitlab.com/ci/testing/code_quality/#golangci-lint).
A simple quickstart is their [CI component](https://gitlab.com/explore/catalog/components/code-quality-oss/codequality-os-scanners-integration), which can be used like this:

```yaml {filename=".gitlab-ci.yml"} 
include:
  - component: $CI_SERVER_FQDN/components/code-quality-oss/codequality-os-scanners-integration/golangci@1.0.1
```

Note that you [can only reference components in the same GitLab instance as your project](https://docs.gitlab.com/ci/components/#use-a-component)

## Other CI

Here are the other ways to install golangci-lint:

{{< cards >}}
  {{< card link="/docs/welcome/install/local/#binaries" title="Bash/Binaries" icon="archive" >}}
  {{< card link="/docs/welcome/install/local/#docker" title="Docker" icon="archive" >}}
{{< /cards >}}


