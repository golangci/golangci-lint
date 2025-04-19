//go:build ignore

// Copyright (c) 2025 golangci-lint <someone@example.com>. // want "invalid license header"
// This file is a part of golangci-lint.

//golangcitest:args -Egolicenser
//golangcitest:config_path testdata/golicenser.yml
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
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
