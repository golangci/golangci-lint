//golangcitest:args -Ethelper
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
	"testing"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func thelperWithHelperAfterAssignment(t *testing.T) { // want "test helper function should start from t.Helper()"
	_ = 0
	t.Helper()
}
