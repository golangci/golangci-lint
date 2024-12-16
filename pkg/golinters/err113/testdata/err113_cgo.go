//golangcitest:args -Eerr113
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
	"os"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func CheckGoerr13Import(e error) bool {
	f, err := os.Create("f.txt")
	if err != nil {
		return err == e // want `do not compare errors directly "err == e", use "errors.Is\(err, e\)" instead`
	}
	f.Close()
	return false
}
