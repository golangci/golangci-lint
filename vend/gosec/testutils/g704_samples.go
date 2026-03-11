package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG704 - SSRF via taint analysis
var SampleCodeG704 = []CodeSample{
	{[]string{`
package main

import (
	"net/http"
)

func handler(r *http.Request) {
	url := r.URL.Query().Get("url")
	http.Get(url)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"net/http"
	"os"
)

func fetchFromEnv() {
	target := os.Getenv("TARGET_URL")
	http.Post(target, "text/plain", nil)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"net/http"
)

func safeRequest() {
	// Safe - hardcoded URL
	http.Get("https://api.example.com/data")
}
`}, 0, gosec.NewConfig()},
}
