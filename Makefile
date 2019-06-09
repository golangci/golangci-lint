.DEFAULT_GOAL = test
.PHONY: FORCE

# Build

build: golangci-lint
clean:
	rm -f golangci-lint test/path
	rm -rf tools
.PHONY: build clean

# Test

test: build
	GL_TEST_RUN=1 ./golangci-lint run -v
	GL_TEST_RUN=1 ./golangci-lint run --fast --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
	GL_TEST_RUN=1 ./golangci-lint run --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
	GL_TEST_RUN=1 go test -v ./...

.PHONY: test

test_race:
	go build -race -o golangci-lint ./cmd/golangci-lint
	GL_TEST_RUN=1 ./golangci-lint run -v --deadline=5m
.PHONY: test_race

test_linters:
	GL_TEST_RUN=1 go test -v ./test -count 1 -run TestSourcesFromTestdataWithIssuesDir/$T
.PHONY: test_linters

# Maintenance

generate: docs/demo.svg README.md install.sh pkg/logutils/log_mock.go vendor
maintainer-clean: clean
	rm -f docs/demo.svg README.md install.sh pkg/logutils/log_mock.go
	rm -rf vendor
.PHONY: generate maintainer-clean

check_generated:
	$(MAKE) --always-make generate
	git diff --exit-code # check no changes
.PHONY: check_generated

release:
	rm -rf dist
	curl -sL https://git.io/goreleaser | bash
.PHONY: release

# Non-PHONY targets (real files)

golangci-lint: FORCE pkg/logutils/log_mock.go
	go build -o $@ ./cmd/golangci-lint

tools/mockgen: go.mod go.sum
	GOBIN=$(CURDIR)/tools GO111MODULE=on go install github.com/golang/mock/mockgen

tools/goimports: go.mod go.sum
	GOBIN=$(CURDIR)/tools GO111MODULE=on go install golang.org/x/tools/cmd/goimports

tools/go.mod:
	@mkdir -p tools
	@rm -f $@
	cd tools && GO111MODULE=on go mod init local-tools

tools/godownloader: Makefile tools/go.mod
	cd tools && GOBIN=$(CURDIR)/tools GO111MODULE=on go get github.com/goreleaser/godownloader@3b90d248ba30307915288f08ab3f2fc2d9f6710c

tools/svg-term:
	@mkdir -p tools
	cd tools && npm install svg-term-cli
	ln -sf node_modules/.bin/svg-term $@

tools/Dracula.itermcolors:
	@mkdir -p tools
	curl -fL -o $@ https://raw.githubusercontent.com/dracula/iterm/master/Dracula.itermcolors

docs/demo.svg: tools/svg-term tools/Dracula.itermcolors
	PATH=$(CURDIR)/tools:$${PATH} svg-term --cast=183662 --out docs/demo.svg --window --width 110 --height 30 --from 2000 --to 20000 --profile ./tools/Dracula.itermcolors --term iterm2

install.sh: tools/godownloader .goreleaser.yml
	# TODO: use when Windows installation will be fixed in the upstream
	#PATH=$(CURDIR)/tools:$${PATH} tools/godownloader .goreleaser.yml | sed '/DO NOT EDIT/s/ on [0-9TZ:-]*//' > $@

README.md: FORCE golangci-lint
	go run ./scripts/gen_readme/main.go

pkg/logutils/log_mock.go: tools/mockgen tools/goimports pkg/logutils/log.go
	@rm -f $@
	PATH=$(CURDIR)/tools:$${PATH} go generate ./...

go.mod: FORCE
	GO111MODULE=on go mod verify
	GO111MODULE=on go mod tidy
go.sum: go.mod

.PHONY: vendor
vendor: go.mod go.sum
	rm -rf vendor
	GO111MODULE=on go mod vendor
