//golangcitest:args -Eerrname
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
	"errors"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var (
	EOF          = errors.New("end of file")
	ErrEndOfFile = errors.New("end of file")
	errEndOfFile = errors.New("end of file")

	EndOfFileError = errors.New("end of file") // want "the sentinel error name `EndOfFileError` should conform to the `ErrXxx` format"
	ErrorEndOfFile = errors.New("end of file") // want "the sentinel error name `ErrorEndOfFile` should conform to the `ErrXxx` format"
	EndOfFileErr   = errors.New("end of file") // want "the sentinel error name `EndOfFileErr` should conform to the `ErrXxx` format"
	endOfFileError = errors.New("end of file") // want "the sentinel error name `endOfFileError` should conform to the `errXxx` format"
	errorEndOfFile = errors.New("end of file") // want "the sentinel error name `errorEndOfFile` should conform to the `errXxx` format"
)
