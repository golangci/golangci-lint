//golangcitest:args -Enilnil
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

func _() (unsafe.Pointer, error) {
	return nil, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}
