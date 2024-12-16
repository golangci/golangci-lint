//golangcitest:args -Evarnamelen
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
	"math"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	x := math.MinInt8 // want "variable name 'x' is too short for the scope of its usage"
	x++
	x++
	x++
	x++
	x++
	x++
	x++
	x++
	x++
}
