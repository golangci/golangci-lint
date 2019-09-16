#! /bin/bash

OUTPUT="$(go list ./... | grep -v vendor | xargs go fmt)"
if [ -n "$OUTPUT" ]; then
    echo "Go code is not formatted, run gofmt on:" >&2
    echo "$OUTPUT" >&2
    false
fi