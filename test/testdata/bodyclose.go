//golangcitest:args -Ebodyclose
package testdata

import (
	"io/ioutil"
	"net/http"
)

func BodycloseNotClosed() {
	resp, _ := http.Get("https://google.com") // want "response body must be closed"
	_, _ = ioutil.ReadAll(resp.Body)
}
