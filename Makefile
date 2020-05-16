.DEFAULT_GOAL = test
.PHONY: FORCE

# enable consistent Go 1.12/1.13 GOPROXY behavior.
export GOPROXY = https://proxy.golang.org

# Build

build: golangci-lint
.PHONY: build

build_race:
	go build -race -o golangci-lint ./cmd/golangci-lint
.PHONY: build_race

clean:
	rm -f golangci-lint
	rm -f test/path
	rm -f tools/Dracula.itermcolors
	rm -f tools/godownloader
	rm -f tools/goreleaser
	rm -f tools/svg-term
	rm -rf tools/node_modules
.PHONY: clean

# Test
test: export GOLANGCI_LINT_INSTALLED = true
test: build
	GL_TEST_RUN=1 ./golangci-lint run -v
	GL_TEST_RUN=1 go test -v -parallel 2 ./...
.PHONY: test

test_race: build_race
	GL_TEST_RUN=1 ./golangci-lint run -v --timeout=5m
.PHONY: test_race

test_linters:
	GL_TEST_RUN=1 go test -v ./test -count 1 -run TestSourcesFromTestdataWithIssuesDir/$T
.PHONY: test_linters

# Maintenance

generate: install.sh assets/github-action-config.json
.PHONY: generate

maintainer-clean: clean
	rm -rf install.sh
.PHONY: maintainer-clean

check_generated:
	$(MAKE) --always-make generate
	git checkout -- go.mod go.sum # can differ between go1.12 and go1.13
	git diff --exit-code # check no changes
.PHONY: check_generated

release: .goreleaser.yml tools/goreleaser
	./tools/goreleaser
.PHONY: release

snapshot: .goreleaser.yml tools/goreleaser
	./tools/goreleaser --snapshot --rm-dist
.PHONY: snapshot

# Non-PHONY targets (real files)

golangci-lint: FORCE
	go build -o $@ ./cmd/golangci-lint

tools/godownloader: export GOFLAGS = -mod=readonly
tools/godownloader: tools/go.mod tools/go.sum
	cd tools && go build github.com/goreleaser/godownloader

tools/goreleaser: export GOFLAGS = -mod=readonly
tools/goreleaser: tools/go.mod tools/go.sum
	cd tools && go build github.com/goreleaser/goreleaser

# TODO: migrate to docs/
tools/svg-term: tools/package.json tools/package-lock.json
	cd tools && npm ci
	ln -sf node_modules/.bin/svg-term $@

# TODO: migrate to docs/
tools/Dracula.itermcolors:
	curl -fL -o $@ https://raw.githubusercontent.com/dracula/iterm/master/Dracula.itermcolors

# TODO: migrate to docs/
assets/demo.svg: tools/svg-term tools/Dracula.itermcolors
	./tools/svg-term --cast=183662 --out assets/demo.svg --window --width 110 --height 30 --from 2000 --to 20000 --profile ./tools/Dracula.itermcolors --term iterm2

install.sh: .goreleaser.yml tools/godownloader
	./tools/godownloader .goreleaser.yml | sed '/DO NOT EDIT/s/ on [0-9TZ:-]*//' > $@

assets/github-action-config.json: FORCE golangci-lint
	go run ./scripts/gen_github_action_config/main.go $@

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod
