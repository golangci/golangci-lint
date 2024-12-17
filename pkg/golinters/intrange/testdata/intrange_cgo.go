//go:build go1.22

//golangcitest:args -Eintrange
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

func CheckIntrange() {
	for i := 0; i < 10; i++ { // want `for loop can be changed to use an integer range \(Go 1\.22\+\)`
	}

	for i := uint8(0); i < math.MaxInt8; i++ { // want `for loop can be changed to use an integer range \(Go 1\.22\+\)`
	}

	for i := 0; i < 10; i += 2 {
	}

	for i := 0; i < 10; i++ {
		i += 1
	}
}
