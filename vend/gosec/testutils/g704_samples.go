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
	{[]string{`
package main

import (
	"context"
	"net/http"
	"time"
)

func GetPublicIP() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://am.i.mullvad.net/ip", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "", nil
}
`}, 0, gosec.NewConfig()},
	// Constant URL string must NOT trigger G704.
	{[]string{`
package main

import (
	"context"
	"net/http"
)

const url = "https://go.dev/"

func main() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		panic(err)
	}
	_, err = new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}
}
`}, 0, gosec.NewConfig()},
	// Sanity check: variable URL from request still fires.
	{[]string{`
package main

import (
	"net/http"
)

func handler(r *http.Request) {
	target := r.URL.Query().Get("url")
	http.Get(target) //nolint:errcheck
}
`}, 1, gosec.NewConfig()},
}
