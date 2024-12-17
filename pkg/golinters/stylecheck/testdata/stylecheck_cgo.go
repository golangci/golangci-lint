//golangcitest:args -Estylecheck
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

func _(x int) {
	switch x {
	case 1:
		return
	default: // want "ST1015: default case should be first or last in switch statement"
		return
	case 2:
		return
	}
}
