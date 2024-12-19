//golangcitest:args -Enilnesserr
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
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func do() error {
	return fmt.Errorf("do error")
}

func do2() error {
	return fmt.Errorf("do2 error")
}

func someCall() error {
	err := do()
	if err != nil {
		return err
	}
	err2 := do2()
	if err2 != nil {
		return err // want `return a nil value error after check error`
	}
	return nil
}

func sameCall2() error {
	err := do()
	if err == nil {
		err2 := do2()
		if err2 != nil {
			return err // want `return a nil value error after check error`
		}
		return nil
	}
	return err

}
