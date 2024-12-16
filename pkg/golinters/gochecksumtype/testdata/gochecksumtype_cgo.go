//golangcitest:args -Egochecksumtype
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
	"log"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

//sumtype:decl
type SumType interface{ isSumType() }

//sumtype:decl
type One struct{} // want "type 'One' is not an interface"

func (One) isSumType() {}

type Two struct{}

func (Two) isSumType() {}

func sumTypeTest() {
	var sum SumType = One{}
	switch sum.(type) { // want "exhaustiveness check failed for sum type.*SumType.*missing cases for Two"
	case One:
	}

	switch sum.(type) { // want "exhaustiveness check failed for sum type.*SumType.*missing cases for Two"
	case One:
	default:
		panic("??")
	}

	switch sum.(type) {
	case *One:
	default:
		log.Println("legit catch all goes here")
	}

	log.Println("??")

	switch sum.(type) {
	case One:
	case Two:
	}
}
