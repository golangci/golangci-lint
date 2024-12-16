//golangcitest:args -Enonamedreturns
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

func _() (i int, err error) { // want `named return "i" with type "int" found`
	defer func() {
		i = math.MaxInt
		err = nil
	}()
	return
}
