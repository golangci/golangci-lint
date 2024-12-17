//golangcitest:args -Egosimple
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
	"log"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _(ss []string) {
	if ss != nil { // want "S1031: unnecessary nil check around range"
		for _, s := range ss {
			log.Printf(s)
		}
	}
}
