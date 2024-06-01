#!/bin/bash -e

# Benchmark between 2 existing versions
# Usage: ./scripts/bench/bench_version.sh gosec v1.58.1 v1.58.2

# ex: gosec
LINTER="$1"

# ex: v1.58.1
VERSION_OLD="$2"
# ex: v1.58.2
VERSION_NEW="$3"

if [ -z "$LINTER" ] || [ -z "$VERSION_OLD" ] || [ -z "$VERSION_NEW" ]; then
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
  rm "./golangci-lint-${VERSION_OLD}"
  rm "./golangci-lint-${VERSION_NEW}"
}

trap cleanBinaries EXIT

## Install

function install() {
  local VERSION=$1

  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./temp-${VERSION}" "${VERSION}"

  mv "temp-${VERSION}/golangci-lint" "./golangci-lint-${VERSION}"
  rm -rf "temp-${VERSION}"
}

## VERSION_OLD

install "${VERSION_OLD}"

## VERSION_NEW

install "${VERSION_NEW}"

## Run

hyperfine \
--prepare "./golangci-lint-${VERSION_OLD} cache clean" "./golangci-lint-${VERSION_OLD} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER}" \
--prepare "./golangci-lint-${VERSION_NEW} cache clean" "./golangci-lint-${VERSION_NEW} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER}"
