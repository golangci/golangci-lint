//golangcitest:args -Enosprintfhostport
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
	ip := "fd00::1"
	_ = fmt.Sprintf("http://%s:1936/healthz", ip) // want "host:port in url should be constructed with net.JoinHostPort and not directly with fmt.Sprintf"
}
