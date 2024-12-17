//golangcitest:args -Egocritic
//golangcitest:config_path testdata/gocritic.yml
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
	"flag"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var _ = *flag.Bool("global1", false, "") // want `flagDeref: immediate deref in \*flag.Bool\(.global1., false, ..\) is most likely an error; consider using flag\.BoolVar`
