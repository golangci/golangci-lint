//golangcitest:args -Ebidichk
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

func _() {
	fmt.Println("LEFT-TO-RIGHT-OVERRIDE: '‭', it is between the single quotes, but it is not visible with a regular editor") // want "found dangerous unicode character sequence LEFT-TO-RIGHT-OVERRIDE"
}
