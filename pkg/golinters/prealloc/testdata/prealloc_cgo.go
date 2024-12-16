//golangcitest:args -Eprealloc
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

func _(source []int) []int {
	var dest []int // want "Consider pre-allocating `dest`"
	for _, v := range source {
		dest = append(dest, v)
	}

	return dest
}
