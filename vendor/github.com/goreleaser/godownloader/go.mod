module github.com/goreleaser/godownloader

go 1.12

// related to an invalid pseudo version in code.gitea.io/gitea v1.10.0-dev.0.20190711052757-a0820e09fbf7
replace github.com/go-macaron/cors => github.com/go-macaron/cors v0.0.0-20190418220122-6fd6a9bfe14e

// related to an invalid pseudo version in contrib.go.opencensus.io/exporter/ocagent@v0.4.2
replace github.com/census-instrumentation/opencensus-proto => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8

require (
	github.com/apex/log v1.1.0
	github.com/client9/codegen v0.0.0-20180316044450-92480ce66a06
	github.com/goreleaser/goreleaser v0.110.0
	github.com/pkg/errors v0.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.2.2
)
