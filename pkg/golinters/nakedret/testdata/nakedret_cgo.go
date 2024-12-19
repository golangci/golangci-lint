//golangcitest:args -Enakedret
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
	"fmt"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() (a int, b string) {
	if a > 0 {
		return // want "naked return in func `_` with 33 lines of code"
	}

	fmt.Println("nakedret")

	if b == "" {
		return 0, "0"
	}

	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...
	// ...

	// len of this function is 33
	return // want "naked return in func `_` with 33 lines of code"
}