# Incident Response Plan

**Last Updated:** 2026-04-11.

This document outlines how the golangci team responds to security incidents, critical bugs,
or operational disruptions that could affect users or the trustworthiness of the project.

Principles:

- Transparency: All incidents and fixes are documented and publicly available.
- Stewardship: Take responsibility for protecting users and the project.
- Protection: Act to minimize harm and provide guidance.

## Scope

This plan applies to:

- The GitHub organization.
- Everything in the following repositories, including code, releases, GitHub workflows, and documentation.
  - [golangci/golangci-lint](https://github.com/golangci/golangci-lint)
  - [golangci/golangci-lint-action](https://github.com/golangci/golangci-lint-action)
  - [golangci/plugin-module-register](https://github.com/golangci/plugin-module-register)
  - [golangci/golines](https://github.com/golangci/golines)
  - [golangci/revgrep](https://github.com/golangci/revgrep)
  - [golangci/misspell](https://github.com/golangci/misspell)
  - [golangci/gofmt](https://github.com/golangci/gofmt)
- Domain names.

## Roles & Contacts

- **Incident Lead:** By default, [@ldez](https://github.com/ldez).
- **Security Contact:** All incidents must be reported exclusively through [GitHub Security Advisories][gsa].

## Detection & Reporting

**All security incidents are initially considered sensitive** and must be reported privately and exclusively through [GitHub Security Advisories][gsa].

Do not disclose incidents through issues, pull requests, or public channels.

## Initial Response

1. **Acknowledge** the report and thank the reporter.
2. **Assess** the severity and validity ([Confidentiality, Integrity, Availability][cia]).
3. **Engage** other maintainers if needed.
4. **Contain** the threat immediately if possible (e.g., revoke credentials, disable workflows).

## Investigation & Mitigation

- **Investigate** the root cause and potential impact.
- **Mitigate**:
    - Patch vulnerabilities.
    - Rotate compromised credentials (tokens/keys).
- **Document** all findings and actions taken.

## Resolution Timeline

Resolution or assessment will typically be provided within **7 business days** from the report date.

## Communication

All communication regarding security incidents must occur exclusively through the [GitHub Security Advisories][gsa].

Once the incident is resolved and a fix is released, we will:

1. Coordinate disclosure timing with the reporter.
2. Publish a public advisory summarizing the incident.
3. Request a CVE identifier if applicable.

## Post-Incident

1. **Review** the incident response and identify lessons learned.
2. **Update** documentation, processes, or automation as needed.
3. **Publish** a public advisory for significant incidents.
4. **Credit** all contributors unless they explicitly request to remain anonymous.

## References

- [Security Policies](./SECURITY.md)
- [GitHub Security Advisories][gsa]
- [CIA (Confidentiality, Integrity, Availability) triad][cia]

[gsa]: https://github.com/golangci/golangci-lint/security/advisories
[cia]: https://www.energy.gov/femp/operational-technology-cybersecurity-energy-systems#cia
