//golangcitest:args -Ecanonicalheader
package testdata

import "net/http"

func canonicalheader() {
	v := http.Header{}

	v.Get("Test-HEader")          // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Set("Test-HEader", "value") // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Add("Test-HEader", "value") // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Del("Test-HEader")          // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Values("Test-HEader")       // want `non-canonical header "Test-HEader", instead use: "Test-Header"`

	v.Set("Test-Header", "value")
	v.Add("Test-Header", "value")
	v.Del("Test-Header")
	v.Values("Test-Header")
}
