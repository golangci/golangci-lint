//golangcitest:args -Enostdglobals
package testdata

import "net/http"

var c = http.Client{}

type inner struct {
	c *http.Client
	DefaultClient string
}

type outer struct {
	f inner
}

var str = outer{inner{c: http.DefaultClient}} // want "should not make use of 'http.DefaultClient'"

var DefaultClient = ""

var _ = http.DefaultClient // want "should not make use of 'http.DefaultClient'"

func f() {

	var _, _ = http.DefaultClient, http.Client{} // want "should not make use of 'http.DefaultClient'"

	_ = http.DefaultClient // want "should not make use of 'http.DefaultClient'"

	_ = http.DefaultTransport // want "should not make use of 'http.DefaultTransport'"

	_ = http.Client{}

	_ = str.f.c
	_ = str.f.DefaultClient
}
