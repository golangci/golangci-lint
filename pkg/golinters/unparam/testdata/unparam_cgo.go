//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Eunparam
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

func unparamUnusedCGO(a, b uint) uint { // want "`unparamUnusedCGO` - `b` is unused"
	a++
	return a
}
