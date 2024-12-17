//golangcitest:args -Enlreturn
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

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func cha() {
	ch := make(chan interface{})
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	select {
	case <-ch:
		return

	case <-ch1:
		{
			a := math.MaxInt
			_ = a
			{
				a := 1
				_ = a
				return // want "return with no blank line before"
			}

			return
		}

		return

	case <-ch2:
		{
			a := 1
			_ = a
			return // want "return with no blank line before"
		}
		return // want "return with no blank line before"
	}
}
