In order for a pull request adding a linter to be reviewed, the linter and the PR must follow some requirements.

- [ ] The CLA must be signed

### Pull Request Description

- [ ] It must have a link to the linter repository.
- [ ] It must provide a short description of the linter.

### Linter

- [ ] It must not be a duplicate of another linter or a rule of a linter. (the team will help to verify that)
- [ ] It must have a valid license (AGPL is not allowed) and the file must contain the required information by the license, ex: author, year, etc.
- [ ] It must use Go <= 1.21
- [ ] The linter repository must have a CI and tests.
- [ ] It must use [`go/analysis`](https://golangci-lint.run/contributing/new-linters/).
- [ ] It must have a valid tag, ex: `v1.0.0`, `v0.1.0`.
- [ ] It must not contain `init()`.
- [ ] It must not contain `panic()`.
- [ ] It must not contain `log.fatal()`, `os.exit()`, or similar.
- [ ] It must not modify the AST.
- [ ] It must not have false positives/negatives. (the team will help to verify that)
- [ ] It must have tests inside golangci-lint.

### The Linter Tests Inside Golangci-lint

- [ ] They must have at least one std lib import.
- [ ] They must have integration tests without configuration (default).
- [ ] They must have integration tests with configuration (if the linter has a configuration).

### `.golangci.next.reference.yml`

- [ ] The file `.golangci.next.reference.yml` must be updated.
- [ ] The file `.golangci.reference.yml` must NOT be edited.
- [ ] The linter must be added to the lists of available linters (alphabetical case-insensitive order).
    - `enable` and `disable` options
- [ ] If the linter has a configuration, the exhaustive configuration of the linter must be added (alphabetical case-insensitive order)
    - The values must be different from the default ones.
    - The default values must be defined in a comment.
    - The option must have a short description.

### Others Requirements

- [ ] The files (tests and linter) inside golangci-lint must have the same name as the linter.
- [ ] The `.golangci.yml` of golangci-lint itself must not be edited and the linter must not be added to this file.
- [ ] The linters must be sorted in the alphabetical order (case-insensitive) in the `lintersdb/builder_linter.go` and `.golangci.next.reference.yml`.
- [ ] The load mode (`WithLoadMode(...)`):
    - if the linter uses `goanalysis.LoadModeSyntax` -> no `WithLoadForGoAnalysis()` in `lintersdb/builder_linter.go`
    - if the linter uses `goanalysis.LoadModeTypesInfo`, it requires `WithLoadForGoAnalysis()` in `lintersdb/builder_linter.go`
- [ ] The version in `WithSince(...)` must be the next minor version (`v1.X.0`) of golangci-lint.
- [ ] `WithURL()` must contain the URL of the repository.

### Recommendations

- [ ] The file `jsonschema/golangci.next.jsonschema.json` should be updated.
- [ ] The file `jsonschema/golangci.jsonschema.json` must NOT be edited.
- [ ] The linter repository should have a readme and linting.
- [ ] The linter should be published as a binary. (useful to diagnose bug origins)
- [ ] The linter repository should have a `.gitignore` (IDE files, binaries, OS files, etc. should not be committed)
- [ ] A tag should never be recreated.

---

The golangci-lint team will edit this comment to check the boxes before and during the review.

The code review will start after the completion of those checkboxes (except for the specific items that the team will help to verify).

The reviews should be addressed as commits (no squash).

If the author of the PR is a member of the golangci-lint team, he should not edit this message.

**This checklist does not imply that we will accept the linter.**
