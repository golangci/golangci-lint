# syntax=docker/dockerfile:1.4
FROM golang:1.22

# related to https://github.com/golangci/golangci-lint/issues/3107
ENV GOROOT /usr/local/go

# Set all directories as safe
RUN git config --global --add safe.directory '*'

COPY golangci-lint /usr/bin/
CMD ["golangci-lint"]
