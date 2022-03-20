# stage 1 building the code
FROM golang:1.18 as builder

ARG VERSION
ARG SHORT_COMMIT
ARG DATE

COPY / /golangci
WORKDIR /golangci
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$VERSION -X main.commit=$SHORT_COMMIT -X main.date=$DATE" -o golangci-lint ./cmd/golangci-lint/main.go

# stage 2
FROM golang:1.18
# don't place it into $GOPATH/bin because Drone mounts $GOPATH as volume
COPY --from=builder /golangci/golangci-lint /usr/bin/
CMD ["golangci-lint"]
