test:
	go install ./cmd/... # needed for govet and golint
	golangci-lint run -v
	golangci-lint run --fast --no-config -v
	golangci-lint run --no-config -v
	go test -v ./...

assets:
	svg-term --cast=183662 --out docs/demo.svg --window --width 110 --height 30 --from 2000 --to 20000 --profile Dracula --term iterm2

readme:
	go run ./scripts/gen_readme/main.go

check_generated:
	make readme && git diff --exit-code # check no changes

release:
	rm -rf dist
	curl -sL https://git.io/goreleaser | bash

.PHONY: test
