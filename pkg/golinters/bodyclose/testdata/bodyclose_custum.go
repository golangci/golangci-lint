//golangcitest:args -Ebodyclose
//golangcitest:config_path testdata/bodyclose_custom.yml
package testdata

import (
	"io"
	"net/http"
)

func consumedWithIOCopy() {
	resp, err := http.Get("http://example.com/") // OK - io.Copy detected
	if err != nil {
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}

func actuallyNotConsumed() {
	resp, err := http.Get("http://example.com/") // want "response body must be closed and consumed"
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
