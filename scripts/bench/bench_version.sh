#!/bin/bash -e

# Benchmark between 2 existing versions
# Usage: ./scripts/bench/bench_version.sh gosec v1.58.1 v1.58.2

# ex: gosec
LINTER_NAME=$1

# ex: v1.58.1
GCIL_VERSION_ONE=$2
# ex: v1.58.2
GCIL_VERSION_TWO=$3

## GCIL_VERSION_ONE

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./temp-${GCIL_VERSION_ONE} ${GCIL_VERSION_ONE}

mv temp-${GCIL_VERSION_ONE}/golangci-lint ./golangci-lint-${GCIL_VERSION_ONE}
rm -rf temp-${GCIL_VERSION_ONE}

## GCIL_VERSION_TWO

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./temp-${GCIL_VERSION_TWO} ${GCIL_VERSION_TWO}

mv temp-${GCIL_VERSION_TWO}/golangci-lint ./golangci-lint-${GCIL_VERSION_TWO}
rm -rf temp-${GCIL_VERSION_TWO}

## Run

hyperfine \
--prepare 'golangci-lint cache clean' "./golangci-lint-${GCIL_VERSION_ONE} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER_NAME}" \
--prepare './golangci-lint cache clean' "./golangci-lint-${GCIL_VERSION_TWO} run --issues-exit-code 0 --print-issued-lines=false --enable-only ${LINTER_NAME}"

## Clean

rm ./golangci-lint-${GCIL_VERSION_ONE}
rm ./golangci-lint-${GCIL_VERSION_TWO}
