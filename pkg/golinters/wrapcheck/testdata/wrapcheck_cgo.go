//golangcitest:args -Ewrapcheck
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
	"encoding/json"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() error {
	_, err := json.Marshal(struct{}{})
	if err != nil {
		return err // want "error returned from external package is unwrapped"
	}

	return nil
}
