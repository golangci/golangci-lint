FROM golang:1.11

COPY golangci-lint /bin/
CMD ["golangci-lint"]
