package cgoexample

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func Example() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	fmt.Printf("bad format %t", cs)
	C.free(unsafe.Pointer(cs))
}
