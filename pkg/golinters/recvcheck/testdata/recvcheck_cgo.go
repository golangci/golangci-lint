//golangcitest:args -Erecvcheck
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

type Bar struct{} // want `the methods of "Bar" use pointer receiver and non-pointer receiver.`

func (b Bar) A() {
	fmt.Println("A")
}

func (b *Bar) B() {
	fmt.Println("B")
}
