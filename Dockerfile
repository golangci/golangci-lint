FROM scratch
COPY golangci-lint /
ENTRYPOINT ["/golangci-lint"]