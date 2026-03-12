# Changelog

## v1.26.1 (2026-03-12)

### Bug Fixes (backported)

- **gosec**: Fix panic on float constants in overflow analyzer ([securego/gosec#1505](https://github.com/securego/gosec/pull/1505))
- **gosec**: Fix panic when scanning multi-module repos from root ([securego/gosec#1504](https://github.com/securego/gosec/pull/1504))
- **gosec**: Fix G407 incorrect detection of fixed IV ([securego/gosec#1509](https://github.com/securego/gosec/pull/1509))
- **gosec**: Fix G115 false positives and false negatives ([securego/gosec#1518](https://github.com/securego/gosec/pull/1518))
- **gosec**: Fix G602 false positive for array element access ([securego/gosec#1499](https://github.com/securego/gosec/pull/1499))
- **gosec**: Fix G704 false positive on const URLs ([securego/gosec#1551](https://github.com/securego/gosec/pull/1551))
- **gosec**: Fix Sonar report schema compliance ([securego/gosec#1507](https://github.com/securego/gosec/pull/1507))
- **gosec**: Fix SARIF output invalid null relationships ([securego/gosec#1569](https://github.com/securego/gosec/pull/1569))
- **unqueryvet**: Fix false positives on constants, prepared statements, and broad name patterns ([v1.5.4](https://github.com/MirrexOne/unqueryvet/releases/tag/v1.5.4))

## v1.26.0 (2026-02-19)

Initial release. Based on `golangci-lint` v2.10.1, targeting Go 1.26.
