//golangcitest:args -Egoprintffuncname
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

func PrintfLikeFuncWithBadName(format string, args ...interface{}) { // want "printf-like formatting function 'PrintfLikeFuncWithBadName' should be named 'PrintfLikeFuncWithBadNamef'"
}
