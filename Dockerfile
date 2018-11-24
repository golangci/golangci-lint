FROM golang:1.11

COPY golangci-lint $GOPATH/bin/
CMD ["golangci-lint"]
