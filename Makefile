test:
	go install ./cmd/...
	golangci-lint run -v
	golangci-lint run --fast --no-config -v
	golangci-lint run --fast --no-config -v
	golangci-lint run --no-config -v
	golangci-lint run --fast --no-config -v ./pkg/testdata/with_issues/typecheck.go
	go test -v -race ./...
