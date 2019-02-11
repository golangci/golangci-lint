FROM golang:1.11.1-alpine3.8 as build
WORKDIR /go/src/github.com/golangci/gosec
COPY . .
RUN apk add -U git make
RUN go get -u github.com/golang/dep/cmd/dep
RUN make

FROM golang:1.11.1-alpine3.8
RUN apk add -U gcc musl-dev
COPY --from=build /go/src/github.com/golangci/gosec/gosec /usr/local/bin/gosec
ENTRYPOINT ["gosec"]
