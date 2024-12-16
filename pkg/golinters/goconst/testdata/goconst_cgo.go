//golangcitest:args -Egoconst
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
	a := "needconst" // want "string `needconst` has 5 occurrences, make it a constant"
	fmt.Print(a)
	b := "needconst"
	fmt.Print(b)
	c := "needconst"
	fmt.Print(c)
}

func _() {
	a := "needconst"
	fmt.Print(a)
	b := "needconst"
	fmt.Print(b)
}
