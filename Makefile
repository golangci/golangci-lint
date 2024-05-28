.DEFAULT_GOAL = test
.PHONY: FORCE

# enable consistent Go 1.12/1.13 GOPROXY behavior.
export GOPROXY = https://proxy.golang.org

BINARY = golangci-lint
ifeq ($(OS),Windows_NT)
	BINARY := $(BINARY).exe
endif

# Build

build: $(BINARY)
.PHONY: build

build_race:
	go build -race -o $(BINARY) ./cmd/golangci-lint
.PHONY: build_race

clean:
	rm -f $(BINARY)
	rm -f test/path
	rm -f tools/Dracula.itermcolors
	rm -f tools/svg-term
	rm -rf tools/node_modules
.PHONY: clean

# Test
test: export GOLANGCI_LINT_INSTALLED = true
test: CGO_ENABLED=1
test: build
	GL_TEST_RUN=1 ./$(BINARY) run -v
	GL_TEST_RUN=1 go test -v -parallel 2 ./...
.PHONY: test

test_race: build_race
	GL_TEST_RUN=1 ./$(BINARY) run -v --timeout=5m
.PHONY: test_race

# ex: T=output.go make test_integration
# the value of `T` is the name of a file from `test/testdata`
test_integration:
	GL_TEST_RUN=1 go test -v ./test -count 1 -run TestSourcesFromTestdata/$T
.PHONY: test_integration

# ex: T=multiple-issues-fix.go make test_integration_fix
# the value of `T` is the name of a file from `test/testdata/fix`
test_integration_fix: build
	GL_TEST_RUN=1 go test -v ./test -count 1 -run TestFix/$T
.PHONY: test_integration_fix

# Maintenance

fast_generate: assets/github-action-config.json
.PHONY: fast_generate

fast_check_generated:
	$(MAKE) --always-make fast_generate
	git checkout -- go.mod go.sum # can differ between go1.16 and go1.17
	git diff --exit-code # check no changes

# Benchmark

# Benchmark with a local version
# LINTER=gosec VERSION=v1.59.0 make bench_local
bench_local: hyperfine
	@:$(call check_defined, LINTER VERSION, 'missing parameter(s)')
	@./scripts/bench/bench_local.sh $(LINTER) $(VERSION)
.PHONY: bench_local

# Benchmark between 2 existing versions
# make bench_version LINTER=gosec VERSION_OLD=v1.58.2 VERSION_NEW=v1.59.0
bench_version: hyperfine
	@:$(call check_defined, LINTER VERSION_OLD VERSION_NEW, 'missing parameter(s)')
	@./scripts/bench/bench_version.sh $(LINTER) $(VERSION_OLD) $(VERSION_NEW)
.PHONY: bench_version

hyperfine:
	@which hyperfine > /dev/null || (echo "Please install hyperfine https://github.com/sharkdp/hyperfine#installation" && exit 1)
.PHONY: hyperfine

# Non-PHONY targets (real files)

$(BINARY): FORCE
	go build -o $@ ./cmd/golangci-lint

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

assets/github-action-config.json: FORCE $(BINARY)
	# go run ./scripts/gen_github_action_config/main.go $@
	cd ./scripts/gen_github_action_config/; go run ./main.go ../../$@

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod

website_copy_jsonschema:
	 go run ./scripts/website/copy_jsonschema/
.PHONY: website_copy_jsonschema

website_expand_templates:
	go run ./scripts/website/expand_templates/
.PHONY: website_expand_templates

website_dump_info:
	go run ./scripts/website/dump_info/
.PHONY: website_dump_info

update_contributors_list:
	cd .github/contributors && npm run all

# Functions

# Check that given variables are set and all have non-empty values,
# die with an error otherwise.
#
# Params:
#   1. Variable name(s) to test.
#   2. (optional) Error message to print.
#
# https://stackoverflow.com/a/10858332/8228109
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))
