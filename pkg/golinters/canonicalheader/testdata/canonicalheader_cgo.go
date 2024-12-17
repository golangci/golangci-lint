//golangcitest:args -Ecanonicalheader
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

func canonicalheader() {
	v := http.Header{}

	v.Get("Test-HEader")          // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Set("Test-HEader", "value") // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Add("Test-HEader", "value") // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Del("Test-HEader")          // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
	v.Values("Test-HEader")       // want `non-canonical header "Test-HEader", instead use: "Test-Header"`

	v.Values("Sec-WebSocket-Accept")

	v.Set("Test-Header", "value")
	v.Add("Test-Header", "value")
	v.Del("Test-Header")
	v.Values("Test-Header")
}
