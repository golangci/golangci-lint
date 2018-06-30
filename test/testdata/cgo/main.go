package cgoexample

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%d\n", s);
}
*/
import "C"

import "unsafe"

func Example() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
