//golangcitest:args -Eforcetypeassert
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
	var a interface{}
	_ = a.(int) // want "type assertion must be checked"

	var b interface{}
	bi := b.(int) // want "type assertion must be checked"
	fmt.Println(bi)
}

func _() {
	var a interface{}
	if ai, ok := a.(int); ok {
		fmt.Println(ai)
	}
}
