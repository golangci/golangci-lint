//golangcitest:args -Ereassign
//golangcitest:config_path testdata/configs/reassign.yml
package testdata

import (
	"io"
	"net/http"
)

func breakIO() {
	http.DefaultClient = nil // want `reassigning variable DefaultClient in other package http`
	io.EOF = nil             // want `reassigning variable EOF in other package io`
}
