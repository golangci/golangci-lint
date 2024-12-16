//golangcitest:args -Enilerr
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

func _() error {
	err := nilErrDo()
	if err == nil {
		return err // want `error is nil \(line 26\) but it returns error`
	}

	return nil
}

func nilErrDo() error {
	return os.ErrNotExist
}
