//golangcitest:args -Elll
//golangcitest:config_path testdata/lll.yml
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

// want +1 "The line is 137 characters long, which exceeds the maximum of 120 characters."
// In my experience, long lines are the lines with comments, not the code. So this is a long comment, a very long comment, yes very long.

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
