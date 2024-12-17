//golangcitest:args -Egovet
//golangcitest:config_path testdata/govet.yml
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
	"io"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _(f io.Reader, buf []byte) (err error) {
	if f != nil {
		_, err := f.Read(buf) // want `shadow: declaration of .err. shadows declaration at line \d+`
		if err != nil {
			return err
		}
	}
	// Use variable to trigger shadowing error
	_ = err
	return
}
