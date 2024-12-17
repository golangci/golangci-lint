//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Ewsl
//golangcitest:config_path testdata/wsl.yml
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

func _() {
	var (
		y = 0
	)
	if y < 1 { // want "if statements should only be cuddled with assignments"
		fmt.Println("tight")
	}
}
