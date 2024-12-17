//golangcitest:args -Eerrcheck
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
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func RetErr() error {
	return nil
}

func MissedErrorCheck() {
	RetErr() // want "Error return value is not checked"
}
