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

	v.Get("Test-HEader")          // want `use "Test-Header" instead of "Test-HEader"`
	v.Set("Test-HEader", "value") // want `use "Test-Header" instead of "Test-HEader"`
	v.Add("Test-HEader", "value") // want `use "Test-Header" instead of "Test-HEader"`
	v.Del("Test-HEader")          // want `use "Test-Header" instead of "Test-HEader"`
	v.Values("Test-HEader")       // want `use "Test-Header" instead of "Test-HEader"`

	v.Values("Sec-WebSocket-Accept")

	v.Set("Test-Header", "value")
	v.Add("Test-Header", "value")
	v.Del("Test-Header")
	v.Values("Test-Header")
}
