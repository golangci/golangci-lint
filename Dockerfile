FROM golang:1.10 AS builder

COPY . $GOPATH/src/github.com/golangci/golangci-lint

RUN make -C $GOPATH/src/github.com/golangci/golangci-lint

FROM golang:1.10

RUN apt-get update && apt-get install -y gcc

COPY --from=builder $GOPATH/bin $GOPATH/bin

ENTRYPOINT ["golangci-lint"]
