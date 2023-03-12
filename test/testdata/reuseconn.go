//golangcitest:args -Ereuseconn
package testdata

import (
	"io"
	"io/ioutil"
	"net/http"
)

func BodyNotDisposedInSingleFunction() {
	resp, _ := http.Get("https://google.com") // want "response body must be disposed properly in a single function read to completion and closed"
	_, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
}

func BodyDisposedInSingleFunction() {
	resp, _ := http.Get("https://google.com")
	disposeResponse(resp)
}

func disposeResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
