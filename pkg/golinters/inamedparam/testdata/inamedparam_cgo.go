//golangcitest:args -Einamedparam
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
	"context"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type tStruct struct {
	a int
}

type Doer interface {
	Do() string
}

type NamedParam interface {
	Void()

	NoArgs() string

	WithName(ctx context.Context, number int, toggle bool, tStruct *tStruct, doer Doer) (bool, error)

	WithoutName(
		context.Context,  // want "interface method WithoutName must have named param for type context.Context"
		int,              // want "interface method WithoutName must have named param for type int"
		bool,             // want "interface method WithoutName must have named param for type bool"
		tStruct,          // want "interface method WithoutName must have named param for type tStruct"
		Doer,             // want "interface method WithoutName must have named param for type Doer"
		struct{ b bool }, // want "interface method WithoutName must have all named params"
	)
}
