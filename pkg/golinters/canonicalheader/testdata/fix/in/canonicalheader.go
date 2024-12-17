//golangcitest:args -Ecanonicalheader
//golangcitest:expected_exitcode 0
package testdata

import "net/http"

func canonicalheader() {
	v := http.Header{}

	v.Get("Test-HEader")
	v.Set("Test-HEader", "value")
	v.Add("Test-HEader", "value")
	v.Del("Test-HEader")
	v.Values("Test-HEader")

	v.Values("Sec-WebSocket-Accept")

	v.Set("Test-Header", "value")
	v.Add("Test-Header", "value")
	v.Del("Test-Header")
	v.Values("Test-Header")
}
