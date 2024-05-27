#!/bin/bash -e

# Benchmark with a local version
# Usage: ./scripts/bench/bench_local.sh gosec v1.59.0

# ex: gosec
LINTER_NAME=$1

# ex: v1.59.0
GCIL_VERSION=$2

## Download version

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./temp-${GCIL_VERSION}/ ${GCIL_VERSION}

mv temp-${GCIL_VERSION}/golangci-lint ./golangci-lint-${GCIL_VERSION}
rm -rf temp-${GCIL_VERSION}

## Build local version

make build

## Run

hyperfine \
--prepare 'golangci-lint cache clean' "./golangci-lint run --print-issued-lines=false --enable-only ${LINTER_NAME}" \
--prepare './golangci-lint cache clean' "./golangci-lint-${GCIL_VERSION} run --print-issued-lines=false --enable-only ${LINTER_NAME}"

## Clean

rm ./golangci-lint-${GCIL_VERSION}
