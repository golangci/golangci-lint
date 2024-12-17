//golangcitest:args -Edupword
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

func duplicateWordInComments() {
	// this line include duplicated word the the // want `Duplicate words \(the\) found`
	fmt.Println("hello")
}
