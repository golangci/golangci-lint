//golangcitest:args -Eexhaustruct
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
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type Exhaustruct struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustruct() {
	// pass
	_ = Exhaustruct{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failPrivate
	_ = Exhaustruct{ // want "testdata.Exhaustruct is missing field c"
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = Exhaustruct{ // want "testdata.Exhaustruct is missing field B"
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = Exhaustruct{ // want "testdata.Exhaustruct is missing fields B, D"
		A: "a",
		c: false,
		E: time.Now(),
	}

}
