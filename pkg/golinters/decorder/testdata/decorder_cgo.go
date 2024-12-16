//golangcitest:args -Edecorder
//golangcitest:config_path testdata/decorder_custom.yml
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
	"math"
	"unsafe"
)

const (
	decoc = math.MaxInt64
	decod = 1
)

var decoa = 1
var decob = 1 // want "multiple \"var\" declarations are not allowed; use parentheses instead"

type decoe int // want "type must not be placed after const"

func decof() {
	const decog = 1
}

func init() {} // want "init func must be the first function in file"

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
