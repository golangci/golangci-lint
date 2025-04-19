//go:build ignore

// TODO(joshuasing): golicenser doesn't currently support cgo, for a few reasons:
//   - source file will be copied to go-build cache
//   - copied source file will contain generated comment, causing it to be excluded
//   - modifying the copied go-build cache file does not result in the actual source file being fixed
//  I would like to eventually support cgo, however for now, if using cgo, it is recommended to add to exclude:
//   - */go-build/*

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
