package main

import (
	"github.com/golangci/golangci-lint/v2/test/testdata_etc/unused_exported/lib"
)

func main() {
	lib.PublicFunc()
}
