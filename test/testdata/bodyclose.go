//args: -Ebodyclose
package testdata

import (
	"io"
	"net/http"
)

func BodycloseNotClosed() {
	resp, _ := http.Get("https://google.com") // ERROR "response body must be closed"
	_, _ = io.ReadAll(resp.Body)
}
