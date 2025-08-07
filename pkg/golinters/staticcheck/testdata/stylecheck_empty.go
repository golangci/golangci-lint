//golangcitest:args -Estaticcheck
//golangcitest:config_path testdata/stylecheck_empty.yml
package testdata

import "net/http"

func _() {
	http.StatusText(200) // want "ST1013: should use constant http.StatusOK instead of numeric literal 200"
	http.StatusText(400) // want "ST1013: should use constant http.StatusBadRequest instead of numeric literal 400"
	http.StatusText(404) // want "ST1013: should use constant http.StatusNotFound instead of numeric literal 404"
	http.StatusText(418) // want "ST1013: should use constant http.StatusTeapot instead of numeric literal 418"
	http.StatusText(500) // want "ST1013: should use constant http.StatusInternalServerError instead of numeric literal 500"
	http.StatusText(503) // want "ST1013: should use constant http.StatusServiceUnavailable instead of numeric literal 503"
	http.StatusText(600)
}
