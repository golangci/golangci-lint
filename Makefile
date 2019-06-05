test: build
	GL_TEST_RUN=1 ./golangci-lint run -v
	GL_TEST_RUN=1 ./golangci-lint run --fast --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
	GL_TEST_RUN=1 ./golangci-lint run --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
	GL_TEST_RUN=1 go test -v ./...

build:
	go build -o golangci-lint ./cmd/golangci-lint

test_race:
	go build -race -o golangci-lint ./cmd/golangci-lint
	GL_TEST_RUN=1 ./golangci-lint run -v --deadline=5m

test_linters:
	GL_TEST_RUN=1 go test -v ./test -count 1 -run TestSourcesFromTestdataWithIssuesDir/$T

assets:
	svg-term --cast=183662 --out docs/demo.svg --window --width 110 --height 30 --from 2000 --to 20000 --profile Dracula --term iterm2

readme:
	go run ./scripts/gen_readme/main.go

gen:
	go generate ./...

check_generated:
	$(MAKE) readme update_deps
	git diff --exit-code # check no changes

release:
	rm -rf dist
	curl -sL https://git.io/goreleaser | bash

update_deps:
	GO111MODULE=on go mod verify
	GO111MODULE=on go mod tidy
	rm -rf vendor
	GO111MODULE=on go mod vendor

.PHONY: test
