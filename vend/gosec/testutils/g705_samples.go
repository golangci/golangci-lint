package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG705 - XSS via taint analysis
var SampleCodeG705 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "<h1>Hello %s</h1>", name)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"net/http"
)

func writeHandler(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	w.Write([]byte(data))
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"net/http"
	"html"
)

func safeHandler(w http.ResponseWriter, r *http.Request) {
	// Safe - escaped output
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "<h1>Hello %s</h1>", html.EscapeString(name))
}
`}, 1, gosec.NewConfig()},
}
