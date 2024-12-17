//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Eineffassign
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
	x := math.MinInt8
	for {
		_ = x
		x = 0 // want "ineffectual assignment to x"
		x = 0
	}
}
