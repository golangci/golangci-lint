//golangcitest:args -Eusestdlibvars
//golangcitest:expected_exitcode 0
package testdata

import "net/http"

func _200() {
	_ = 200
}

func _200_1() {
	var w http.ResponseWriter
	w.WriteHeader(200) // want `"200" can be replaced by http.StatusOK`
}
