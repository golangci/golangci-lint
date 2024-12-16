//golangcitest:args -Efunlen
//golangcitest:config_path testdata/funlen.yml
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

func _() { // want `Function '_' is too long \(22 > 20\)`
	t := struct {
		A string
		B string
		C string
		D string
		E string
		F string
		G string
		H string
		I string
	}{
		`a`,
		`b`,
		`c`,
		`d`,
		`e`,
		`f`,
		`g`,
		`h`,
		`i`,
	}
	_ = t
}
