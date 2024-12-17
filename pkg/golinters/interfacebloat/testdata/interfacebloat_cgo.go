//golangcitest:args -Einterfacebloat
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
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type _ interface { // want "the interface has more than 10 methods: 11"
	a01() time.Duration
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
}
