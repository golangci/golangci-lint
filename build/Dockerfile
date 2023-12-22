# stage 1 building the code
FROM golang:1.21 as builder

ARG VERSION
ARG SHORT_COMMIT
ARG DATE

COPY / /golangci
WORKDIR /golangci
RUN APP_VERSION=${VERSION#v} \
    CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w -X main.version=$APP_VERSION -X main.commit=$SHORT_COMMIT -X main.date=$DATE" -o golangci-lint ./cmd/golangci-lint/main.go

# stage 2
FROM golang:1.21
# related to https://github.com/golangci/golangci-lint/issues/3107
ENV GOROOT /usr/local/go
# Set all directories as safe
RUN git config --global --add safe.directory '*'
# don't place it into $GOPATH/bin because Drone mounts $GOPATH as volume
COPY --from=builder /golangci/golangci-lint /usr/bin/
CMD ["golangci-lint"]
