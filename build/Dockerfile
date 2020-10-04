FROM golang:1.15

# don't place it into $GOPATH/bin because Drone mounts $GOPATH as volume
COPY golangci-lint /usr/bin/
CMD ["golangci-lint"]
