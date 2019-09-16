SHELL       = /bin/bash -euo pipefail
PKGS        = go list ./... | grep -v vendor
GO111MODULE = on
GOFLAGS     = -mod=vendor
TIMEOUT     = 1s


.DEFAULT_GOAL = test-with-coverage


.PHONY: deps
deps:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: update
update:
	@go get -mod= -u


.PHONY: format
format:
	@goimports -local $(dirname $(go list -m)) -ungroup -w .

.PHONY: generate
generate:
	@go generate ./...

.PHONY: refresh
refresh: generate format


.PHONY: test
test:
	@go test -race -timeout $(TIMEOUT) ./...

.PHONY: test-with-coverage
test-with-coverage:
	@go test -cover -timeout $(TIMEOUT) ./... | column -t | sort -r

.PHONY: test-with-coverage-profile
test-with-coverage-profile:
	@go test -cover -covermode count -coverprofile c.out -timeout $(TIMEOUT) ./...


.PHONY: sync
sync:
	@git stash && git pull --rebase && git stash pop || true

.PHONY: upgrade
upgrade: sync update deps refresh test-with-coverage
