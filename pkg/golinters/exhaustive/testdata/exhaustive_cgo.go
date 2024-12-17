//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Eexhaustive
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

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func processDirection(d Direction) {
	switch d { // want "missing cases in switch of type testdata.Direction: testdata.East, testdata.West"
	case North, South:
	}
}
