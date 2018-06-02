test:
	go install ./cmd/...
	golangci-lint run -v
	golangci-lint run --fast --no-config -v
	golangci-lint run --no-config -v
	golangci-lint run --fast --no-config -v ./test/testdata/typecheck.go
	go test -v -race ./...

assets:
	svg-term --cast=183662 --out docs/demo.svg --window --width 110 --height 30 --from 2000 --to 20000 --profile Dracula --term iterm2

readme:
	go run ./scripts/gen_readme/main.go

check_generated:
	make readme && git diff --exit-code # check no changes

.PHONY: test