//golangcitest:args -Ecanonicalheader
//golangcitest:config_path testdata/canonicalheader_custom.yml
package testdata

import "net/http"

func _() {
	v := http.Header{}

	v.Get("Test-HEader")      // want `use "Test-Header" instead of "Test-HEader"`
	v.Get("WWW-Authenticate") // want `Www-Authenticate" instead of "WWW-Authenticate"`
	v.Get("XXX-XXX")
}
