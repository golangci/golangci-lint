//golangcitest:args -Epredeclared
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
	var real int // want "real: shadows predeclared identifier"
	a := A{}
	copy := Clone(a) // want "copy: shadows predeclared identifier"

	// suppress any "declared but not used" errors
	_ = real
	_ = a
	_ = copy
}

type A struct {
	true bool
	foo  int
}

func Clone(a A) A {
	return A{
		true: a.true,
		foo:  a.foo,
	}
}
