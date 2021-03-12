module github.com/golangci/golangci-lint/tools

go 1.13

require github.com/goreleaser/goreleaser v0.155.0

// https://github.com/mattn/go-shellwords/pull/39
replace github.com/mattn/go-shellwords => github.com/caarlos0/go-shellwords v1.0.11
