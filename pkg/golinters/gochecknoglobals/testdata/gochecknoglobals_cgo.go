//golangcitest:args -Egochecknoglobals
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
	"fmt"
	"regexp"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var noGlobalsVar int // want "noGlobalsVar is a global variable"
var ErrSomeType = errors.New("test that global errors aren't warned")

var (
	OnlyDigits  = regexp.MustCompile(`^\d+$`)
	BadNamedErr = errors.New("this is bad") // want "BadNamedErr is a global variable"
)

func NoGlobals() {
	fmt.Print(noGlobalsVar)
}
