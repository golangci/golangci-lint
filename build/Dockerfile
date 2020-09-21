# stage 1 building the code
FROM golang:1.15 as builder

COPY / /golangci
WORKDIR /golangci
RUN go build -o golangci-lint ./cmd/golangci-lint/main.go

# stage 2
FROM golang:1.15
# don't place it into $GOPATH/bin because Drone mounts $GOPATH as volume
COPY --from=builder /golangci/golangci-lint /usr/bin/
CMD ["golangci-lint"]
