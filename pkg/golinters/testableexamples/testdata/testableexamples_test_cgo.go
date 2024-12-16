//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Etestableexamples
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

func Example_good() {
	fmt.Println("hello")
	// Output: hello
}

func Example_goodEmptyOutput() {
	fmt.Println("")
	// Output:
}

func Example_bad() { // want `^missing output for example, go test can't validate it$`
	fmt.Println("hello")
}

//nolint:testableexamples
func Example_nolint() {
	fmt.Println("hello")
}
