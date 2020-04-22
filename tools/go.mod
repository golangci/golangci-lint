module github.com/golangci/golangci-lint/tools

go 1.12

require (
	cloud.google.com/go v0.52.0 // indirect
	cloud.google.com/go/storage v1.5.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-sdk-for-go v39.0.0+incompatible // indirect
	github.com/Azure/go-autorest v13.3.3+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.9.5 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/apex/log v1.1.2 // indirect
	github.com/aws/aws-sdk-go v1.28.13 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/wire v0.4.0 // indirect
	github.com/goreleaser/godownloader v0.1.0
	github.com/goreleaser/goreleaser v0.132.0
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/kamilsk/retry/v4 v4.4.2 // indirect
	github.com/mattn/go-ieproxy v0.0.0-20200203040449-2dbc853185d9 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/crypto v0.0.0-20200208060501-ecb85df21340 // indirect
	golang.org/x/exp v0.0.0-20200207192155-f17229e696bd // indirect
	golang.org/x/lint v0.0.0-20200130185559-910be7a94367 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/tools v0.0.0-20200207224406-61798d64f025 // indirect
	google.golang.org/api v0.17.0 // indirect
	google.golang.org/genproto v0.0.0-20200210034751-acff78025515 // indirect
	google.golang.org/grpc v1.27.1 // indirect
)

// Fix godownloader/goreleaser deps (ambiguous imports/invalid pseudo-version)
// https://github.com/goreleaser/goreleaser/issues/1145
replace github.com/go-macaron/cors => github.com/go-macaron/cors v0.0.0-20190418220122-6fd6a9bfe14e
