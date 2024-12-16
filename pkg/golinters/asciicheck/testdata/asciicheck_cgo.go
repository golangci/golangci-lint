//golangcitest:args -Easciicheck
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
	"fmt"
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type AsciicheckTеstStruct struct { // want `identifier "AsciicheckTеstStruct" contain non-ASCII character: U\+0435 'е'`
	Date time.Time
}

type AsciicheckField struct{}

type AsciicheckJustStruct struct {
	Tеst AsciicheckField // want `identifier "Tеst" contain non-ASCII character: U\+0435 'е'`
}

func AsciicheckTеstFunc() { // want `identifier "AsciicheckTеstFunc" contain non-ASCII character: U\+0435 'е'`
	var tеstVar int // want `identifier "tеstVar" contain non-ASCII character: U\+0435 'е'`
	tеstVar = 0
	fmt.Println(tеstVar)
}
