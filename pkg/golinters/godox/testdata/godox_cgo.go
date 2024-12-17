//golangcitest:args -Egodox
//golangcitest:config_path testdata/godox.yml
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

func todoLeftInCode() {
	// TODO implement me // want `Line contains FIXME/TODO: "TODO implement me`
	//TODO no space // want `Line contains FIXME/TODO: "TODO no space`
	// TODO(author): 123 // want `Line contains FIXME/TODO: "TODO\(author\): 123`
	//TODO(author): 123 // want `Line contains FIXME/TODO: "TODO\(author\): 123`
	//TODO(author) 456 // want `Line contains FIXME/TODO: "TODO\(author\) 456`
	// TODO: qwerty // want `Line contains FIXME/TODO: "TODO: qwerty`
	// todo 789 // want `Line contains FIXME/TODO: "todo 789`
}
