//golangcitest:args -Ebodyclose
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
	"io/ioutil"
	"net/http"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func BodycloseNotClosed() {
	resp, _ := http.Get("https://google.com") // want "response body must be closed"
	_, _ = ioutil.ReadAll(resp.Body)
}
