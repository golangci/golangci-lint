FROM golang:1.11

RUN apt-get update && apt-get install -y gcc

COPY golangci-lint $GOPATH/bin/
ENTRYPOINT ["golangci-lint"]
