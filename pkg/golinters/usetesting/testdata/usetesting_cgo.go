//golangcitest:args -Eusetesting
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
	"testing"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "") // want `os\.MkdirTemp\(\) could be replaced by t\.TempDir\(\) in .+`
}
