module github.com/golangci/golangci-lint/tools

go 1.12

require (
	github.com/goreleaser/godownloader v0.0.0-20191002112816-e64d0375716b
	github.com/goreleaser/goreleaser v0.118.2
)

// Fix godownloader/goreleaser deps (ambiguous imports/invalid pseudo-version)
// https://github.com/goreleaser/goreleaser/issues/1145
replace github.com/go-macaron/cors => github.com/go-macaron/cors v0.0.0-20190418220122-6fd6a9bfe14e
