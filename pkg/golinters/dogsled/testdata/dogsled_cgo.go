//golangcitest:args -Edogsled
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

func _() {
	_ = ret1()
	_, _ = ret2()
	_, _, _ = ret3()    // want "declaration has 3 blank identifiers"
	_, _, _, _ = ret4() // want "declaration has 4 blank identifiers"
}

func ret1() (a int) {
	return 1
}

func ret2() (a, b int) {
	return 1, 2
}

func ret3() (a, b, c int) {
	return 1, 2, 3
}

func ret4() (a, b, c, d int) {
	return 1, 2, 3, 4
}
