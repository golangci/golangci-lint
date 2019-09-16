SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export PATH := ./bin:$(PATH)
export GO111MODULE := on
export GOPROXY := https://gocenter.io

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod download
.PHONY: setup


pull_test_imgs:
	grep FROM ./acceptance/testdata/*.dockerfile | cut -f2 -d' ' | sort | uniq | while read -r img; do docker pull "$$img"; done
.PHONY: pull_test_imgs

test: pull_test_imgs
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m
.PHONY: test

cover: test
	go tool cover -html=coverage.out
.PHONY: cover

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

lint: check
	./bin/golangci-lint run --enable-all ./...
.PHONY: check

ci: build lint test
.PHONY: ci

build:
	go build -o nfpm ./cmd/nfpm/main.go
.PHONY: build

todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=bin \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build
