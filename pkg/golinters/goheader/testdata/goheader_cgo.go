//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

/*MY TITLE!*/ // want `Expected:TITLE\., Actual: TITLE!`

//golangcitest:args -Egoheader
//golangcitest:config_path testdata/goheader.yml
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
