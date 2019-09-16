SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export PATH := ./bin:$(PATH)
export GO111MODULE := on
export GOPROXY := https://gocenter.io

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/gohugoio/hugo.sh | sh
	curl -L https://git.io/misspell | sh
	go mod tidy
.PHONY: setup

# Run all the tests
test:
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

# Run all the linters
lint:
	# TODO: fix tests issues
	# TODO: fix lll issues
	# TODO: fix funlen issues
	./bin/golangci-lint run --tests=false --enable-all --disable=lll --disable funlen ./...
	./bin/misspell -error **/*
.PHONY: lint

# Clean go.mod
go-mod-tidy:
	@go mod tidy -v
	@git diff HEAD
	@git diff-index --quiet HEAD
.PHONY: go-mod-tidy

# Run all the tests and code checks
ci: build test lint go-mod-tidy
.PHONY: ci

# Build a beta version of goreleaser
build:
	go build
.PHONY: build

# Generate the static documentation
static:
	@hugo --enableGitInfo --source www
.PHONY: static

imgs:
	wget -O www/static/card.png "https://og.caarlos0.dev/**GoReleaser**%20%7C%20Deliver%20Go%20binaries%20as%20fast%20and%20easily%20as%20possible.png?theme=light&md=1&fontSize=80px&images=https://github.com/goreleaser.png"
	wget -O www/static/avatar.png https://github.com/goreleaser.png
	convert www/static/avatar.png -define icon:auto-resize=64,48,32,16 www/static/favicon.ico
	convert www/static/avatar.png -resize x120 www/static/apple-touch-icon.png
.PHONY: imgs

serve: imgs
	@hugo server --enableGitInfo --watch --source www
.PHONY: serve

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build
