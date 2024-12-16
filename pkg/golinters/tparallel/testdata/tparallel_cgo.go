//golangcitest:args -Etparallel
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

func TestSubtests(t *testing.T) { // want "TestSubtests's subtests should call t.Parallel"
	t.Parallel()

	t.Run("", func(t *testing.T) {
	})
}
