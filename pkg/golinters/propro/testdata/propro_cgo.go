//golangcitest:args -Epropro
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

type CgoEntity struct {
	IntField int
}

func (s *CgoEntity) SetProtectedField(value int) {
	s.IntField = value
}

func CgoFunc1() {
	e := &CgoEntity{}
	e.SetProtectedField(1)
	e.IntField = 10 // want "assignment to exported field CgoEntity.IntField is forbidden outside its methods"
}
