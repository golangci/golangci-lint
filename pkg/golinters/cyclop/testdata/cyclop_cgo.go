//golangcitest:args -Ecyclop
//golangcitest:config_path testdata/cyclop.yml
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

func cyclopComplexFunc(s string) { // want "calculated cyclomatic complexity for function cyclopComplexFunc is 22, max is 15"
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
}
