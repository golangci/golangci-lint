//golangcitest:args -Eusestdlibvars
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"net/http"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _200() {
	_ = 200
}

func _200_1() {
	var w http.ResponseWriter
	w.WriteHeader(200) // want `"200" can be replaced by http.StatusOK`
}
