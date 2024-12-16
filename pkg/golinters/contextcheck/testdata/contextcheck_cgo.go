//golangcitest:args -Econtextcheck
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

func contextcheckCase1(ctx context.Context) {
	funcWithoutCtx() // want "Function `funcWithoutCtx` should pass the context parameter"
}

func funcWithCtx(ctx context.Context) {}

func funcWithoutCtx() {
	funcWithCtx(context.TODO())
}
