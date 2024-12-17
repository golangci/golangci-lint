//golangcitest:args -Egocheckcompilerdirectives
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
	_ "embed"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

// Okay cases:

//go:generate echo hello world

//go:embed
var Value string

//go:

// Problematic cases:

// go:embed // want "compiler directive contains space: // go:embed"

//    go:embed // want "compiler directive contains space: //    go:embed"

//go:genrate // want "compiler directive unrecognized: //go:genrate"
