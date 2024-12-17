//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Egosec
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
	"crypto/md5"
	"log"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func Gosec() {
	h := md5.New() // want "G401: Use of weak cryptographic primitive"
	log.Print(h)
}
