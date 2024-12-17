//golangcitest:args -Ereassign
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
	"io"
	"net/http"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	http.DefaultClient = nil
	http.DefaultTransport = nil
	io.EOF = nil // want `reassigning variable EOF in other package io`
}
