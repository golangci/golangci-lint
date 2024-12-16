//golangcitest:args -Emakezero
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

func _() []int {
	x := make([]int, math.MaxInt8)
	return append(x, 1) // want "append to slice `x` with non-zero initialized length"
}
