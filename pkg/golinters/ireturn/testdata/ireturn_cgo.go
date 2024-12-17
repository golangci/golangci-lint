//golangcitest:args -Eireturn
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

type (
	IreturnDoer interface{ Do() }
	ireturnDoer struct{}
)

func _() IreturnDoer       { return new(ireturnDoer) } // want `_ returns interface \(command-line-arguments.IreturnDoer\)`
func (d *ireturnDoer) Do() { /*...*/ }

func _() *ireturnDoer { return new(ireturnDoer) }
