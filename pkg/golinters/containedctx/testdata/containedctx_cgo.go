//golangcitest:args -Econtainedctx
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
	"context"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type ok struct {
	i int
	s string
}

type ng struct {
	ctx context.Context // want "found a struct that contains a context.Context field"
}

type empty struct{}
