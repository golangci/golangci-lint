---
title: Roadmap
weight: 4
aliases:
  - /product/roadmap/
---

## 💡 Feature Requests

Please file an issue to suggest new features. Vote on feature requests by adding a 👍. This helps maintainers prioritize what to work on.

[See Feature Requests](https://github.com/golangci/golangci-lint/issues?utf8=✓&q=is%3Aissue+is%3Aopen+sort%3Areactions-%2B1-desc+label%3Aenhancement)

## 🐛 Bugs

Please file an issue for bugs, missing documentation or unexpected behavior.

[See Bugs](https://github.com/golangci/golangci-lint/issues?utf8=✓&q=is%3Aissue+is%3Aopen+label%3A%22bug%22+sort%3Acreated-desc)

## Versioning Policy

golangci-lint follows [semantic versioning](https://semver.org). However, due to the nature of golangci-lint as a code quality tool,
it's not always clear when a minor or major version bump occurs.
To help clarify this for everyone, we've defined the following semantic versioning policy:

- Patch release (intended to not break your lint build)
  - A patch version update in a specific linter that results in golangci-lint reporting fewer errors.
  - A bug fix to the CLI or core (packages loading, runner, postprocessors, etc).
  - Improvements to documentation.
  - Non-user-facing changes such as refactoring code, adding, deleting, or modifying tests, and increasing test coverage.
  - Re-releasing after a failed release (i.e., publishing a release that doesn't work for anyone).
- Minor release (might break your lint build because of newly found issues)
  - A major or minor version update of a specific linter that results in golangci-lint reporting more errors.
  - A new linter is added.
  - An existing configuration option or linter is deprecated.
  - A new CLI command is created.
  - Backward incompatible change of configuration.
- Major release (likely to break your lint build)
  - Backward incompatible change of configuration with huge impact.

According to our policy, any minor update may report more errors than the previous release (ex: from a bug fix).
As such, we recommend using the fixed minor version and fixed or the latest patch version to guarantee the results of your builds.

For example, in our [GitHub Action](https://github.com/golangci/golangci-lint-action) we require users to explicitly set the minor version of golangci-lint
and we always use the latest patch version.

## Linter Deprecation Cycle

A linter can be deprecated for various reasons, e.g. the linter stops working with a newer version of Go or the author has abandoned its linter.

The deprecation of a linter will follow 3 phases:

1. **Display of a warning message**: The linter can still be used (unless it's completely non-functional),
  but it's recommended to remove it from your configuration.
2. **Display of an error message**: At this point, you should remove the linter.
  The original implementation is replaced by a placeholder that does nothing.
  The linter is NOT enabled when using `default: all` and should be removed from the `disable` option.
3. **Removal of the linter** from golangci-lint.

Each phase corresponds to a minor version:

- v1.0.0 -> warning message
- v1.1.0 -> error message
- v1.2.0 -> linter removed

We will provide clear information about those changes on different supports: changelog, logs, social network, etc.

We consider the removal of a linter as non-breaking changes for golangci-lint itself.
No major version will be created when a linter is removed.

## Future Plans

1. Upstream all changes of forked linters.
2. Make it easy to write own linter/checker: it should take a minimum code, have perfect documentation, debugging and testing tooling.
3. Speed up SSA loading: on-disk cache and existing code profiling-optimizing.
4. Analyze (don't only filter) only new code: analyze only changed files and dependencies, make incremental analysis, caches.
5. Smart new issues detector: don't print existing issues on changed lines.
6. Minimize false-positives by fixing linters and improving testing tooling.
7. Documentation for every issue type.
