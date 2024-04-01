//golangcitest:args -Ereassign
//golangcitest:config_path testdata/reassign_patterns.yml
package testdata

import (
	"io"
	"net/http"
)

func reassignTestPatterns() {
	http.DefaultClient = nil    // want `reassigning variable DefaultClient in other package http`
	http.DefaultTransport = nil // want `reassigning variable DefaultTransport in other package http`
	io.EOF = nil
}
