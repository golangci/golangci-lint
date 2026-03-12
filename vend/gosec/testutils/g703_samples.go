package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG703 - Path traversal via taint analysis
var SampleCodeG703 = []CodeSample{
	{[]string{`
package main

import (
	"net/http"
	"os"
)

func handler(r *http.Request) {
	path := r.URL.Query().Get("file")
	os.Open(path)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"net/http"
	"os"
)

func writeHandler(r *http.Request) {
	filename := r.FormValue("name")
	os.WriteFile(filename, []byte("data"), 0644)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"os"
)

func safeOpen() {
	// Safe - no user input
	os.Open("/var/log/app.log")
}
`}, 0, gosec.NewConfig()},
}
