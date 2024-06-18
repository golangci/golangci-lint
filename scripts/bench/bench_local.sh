#!/bin/bash -e

# Benchmark with a local version
# Usage: ./scripts/bench/bench_local.sh gosec v1.59.0

# ex: gosec
LINTER=$1

# ex: v1.59.0
VERSION=$2


if [ -z "$LINTER" ] || [ -z "$VERSION" ]; then
  cat <<-EOF
Missing required arguments!

Usage:   $0 <linter> <old version> <new version>
Example: $0 gosec v1.58.1 v1.58.2
EOF

  exit 1
fi

## Clean

function cleanBinaries() {
  echo "Clean binaries"
  rm "./golangci-lint-${VERSION}"
  rm ./golangci-lint
}

trap cleanBinaries EXIT

## Download version

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./temp-${VERSION}" "${VERSION}"

mv "temp-${VERSION}/golangci-lint" "./golangci-lint-${VERSION}"
rm -rf "temp-${VERSION}"

## Build local version
## use `go build` to set ldflags (it reduces some performance differences with binaries created by goreleaser)

go build -trimpath -ldflags '-s -w' -o golangci-lint ./cmd/golangci-lint

## Run

hyperfine --warmup 1 \
-n 'local' --prepare './golangci-lint cache clean' "./golangci-lint run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER}" \
-n "${VERSION}" --prepare "./golangci-lint-${VERSION} cache clean" "./golangci-lint-${VERSION} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER}"
