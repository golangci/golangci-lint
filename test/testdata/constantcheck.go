//golangcitest:args -Econstantcheck
package testdata

import (
	"net/http"
)

func checkGet() {
	var w http.ResponseWriter
	http.Error(w, "GET", 200) // want `GET literal contains in constant with name MethodGet`
}
