package main

import (
	"github.com/anduril/golangci-lint/test/testdata_etc/unused_exported/lib"
)

func main() {
	lib.PublicFunc()
}
