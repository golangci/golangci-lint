# syntax=docker/dockerfile:1.4
FROM golang:1.22-alpine

# related to https://github.com/golangci/golangci-lint/issues/3107
ENV GOROOT /usr/local/go

# gcc is required to support cgo;
# git and mercurial are needed most times for go get`, etc.
# See https://github.com/docker-library/golang/issues/80
RUN apk --no-cache add gcc musl-dev git mercurial

# Set all directories as safe
RUN git config --global --add safe.directory '*'

COPY golangci-lint /usr/bin/
CMD ["golangci-lint"]
