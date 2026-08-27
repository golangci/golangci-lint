//golangcitest:args -Ecanonicalheader
package testdata

import "net/http"

func canonicalheader() {
	v := http.Header{}

	v.Get("Test-HEader")          // want `use "Test-Header" instead of "Test-HEader"`
	v.Set("Test-HEader", "value") // want `use "Test-Header" instead of "Test-HEader"`
	v.Add("Test-HEader", "value") // want `use "Test-Header" instead of "Test-HEader"`
	v.Del("Test-HEader")          // want `use "Test-Header" instead of "Test-HEader"`
	v.Values("Test-HEader")       // want `use "Test-Header" instead of "Test-HEader"`

	v.Values("Sec-WebSocket-Accept")

	v.Set("Test-Header", "value")
	v.Add("Test-Header", "value")
	v.Del("Test-Header")
	v.Values("Test-Header")
}
