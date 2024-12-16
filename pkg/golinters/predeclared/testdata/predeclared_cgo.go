//golangcitest:args -Epredeclared
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

func _() {
	var real int // want "variable real has same name as predeclared identifier"
	a := A{}
	copy := Clone(a) // want "variable copy has same name as predeclared identifier"

	// suppress any "declared but not used" errors
	_ = real
	_ = a
	_ = copy
}
