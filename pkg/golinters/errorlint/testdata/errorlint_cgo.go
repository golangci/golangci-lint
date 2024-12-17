//golangcitest:args -Eerrorlint
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
	"log"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var errLintFoo = errors.New("foo")

type errLintBar struct{}

func (*errLintBar) Error() string {
	return "bar"
}

func errorLintAll() {
	err := func() error { return nil }()
	if err == errLintFoo { // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		log.Println("errCompare")
	}

	err = errors.New("oops")
	fmt.Errorf("error: %v", err) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"

	switch err.(type) { // want "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *errLintBar:
		log.Println("errLintBar")
	}
}
