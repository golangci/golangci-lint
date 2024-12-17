//golangcitest:args -Easasalint
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

func getArgsLength(args ...interface{}) int {
	// this line will not report as error
	fmt.Println(args)
	return len(args)
}

func checkArgsLength(args ...interface{}) int {
	return getArgsLength(args) // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.interface\{\}\)`
}

func someCall() {
	var a = []interface{}{1, 2, 3}
	fmt.Println(checkArgsLength(a...) == getArgsLength(a)) // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.interface\{\}\)`
	fmt.Println(checkArgsLength(a...) == getArgsLength(a...))
}
