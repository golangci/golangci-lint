module github.com/golangci/golangci-lint/tools

go 1.12

require (
	github.com/goreleaser/godownloader v0.1.0
	github.com/goreleaser/goreleaser v0.134.0
)

// Fix invalid pseudo-version: revision is longer than canonical (6fd6a9bfe14e)
replace github.com/go-macaron/cors => github.com/go-macaron/cors v0.0.0-20190418220122-6fd6a9bfe14e
