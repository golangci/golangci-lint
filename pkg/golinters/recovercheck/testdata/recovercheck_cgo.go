//golangcitest:args -Erecovercheck
package testdata

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import (
	"log"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from CGO\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

// UnsafeCGOGoroutine demonstrates unsafe goroutine with CGO calls
func UnsafeCGOGoroutine() {
	go func() { // want "goroutine created without panic recovery"
		cs := C.CString("This will crash")
		C.myprint(cs)
		C.free(unsafe.Pointer(cs))
		panic("CGO panic without recovery")
	}()
}

// SafeCGOGoroutine demonstrates safe goroutine with CGO calls and panic recovery
func SafeCGOGoroutine() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered CGO panic: %v", r)
			}
		}()
		cs := C.CString("This is safe")
		C.myprint(cs)
		C.free(unsafe.Pointer(cs))
		panic("CGO panic with recovery")
	}()
}
