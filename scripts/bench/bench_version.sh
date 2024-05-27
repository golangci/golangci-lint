#!/bin/bash -e

# Benchmark between 2 existing versions
# Usage: ./scripts/bench/bench_version.sh gosec v1.58.1 v1.58.2

# ex: gosec
LINTER="$1"

# ex: v1.58.1
VERSION_OLD="$2"
# ex: v1.58.2
VERSION_NEW="$3"

## Clean

function cleanBinaries() {
  echo "Clean binaries"
  rm ./golangci-lint-${VERSION_OLD}
  rm ./golangci-lint-${VERSION_NEW}
}

trap cleanBinaries EXIT

## VERSION_OLD

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./temp-${VERSION_OLD} ${VERSION_OLD}

mv temp-${VERSION_OLD}/golangci-lint ./golangci-lint-${VERSION_OLD}
rm -rf temp-${VERSION_OLD}

## VERSION_NEW

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./temp-${VERSION_NEW} ${VERSION_NEW}

mv temp-${VERSION_NEW}/golangci-lint ./golangci-lint-${VERSION_NEW}
rm -rf temp-${VERSION_NEW}

## Run

hyperfine \
--prepare 'golangci-lint cache clean' "./golangci-lint-${VERSION_OLD} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER}" \
--prepare './golangci-lint cache clean' "./golangci-lint-${VERSION_NEW} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER}"
