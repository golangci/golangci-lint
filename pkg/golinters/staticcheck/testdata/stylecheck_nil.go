//golangcitest:args -Estaticcheck
//golangcitest:config_path testdata/stylecheck_nil.yml
package testdata

import "net/http"

func _() {
	http.StatusText(200)
	http.StatusText(400)
	http.StatusText(404)
	http.StatusText(418) // want "ST1013: should use constant http.StatusTeapot instead of numeric literal 418"
	http.StatusText(500)
	http.StatusText(503) // want "ST1013: should use constant http.StatusServiceUnavailable instead of numeric literal 503"
	http.StatusText(600)
}
