//golangcitest:args -Egoimports
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
	"unsafe" // want "File is not properly formatted"
	"github.com/golangci/golangci-lint/pkg/config"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func Bar() {
	fmt.Print("x")
	_ = config.Config{}
}
