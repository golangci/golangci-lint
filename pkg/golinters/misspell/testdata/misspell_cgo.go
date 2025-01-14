//golangcitest:args -Emisspell
//golangcitest:config_path testdata/misspell.yml
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

func Misspell() {
	// comment with incorrect spelling: occured // want "`occured` is a misspelling of `occurred`"
}

// the word langauge should be ignored here: it's set in config
// the word Dialogue should be ignored here: it's set in config

func _() error {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))

	return fmt.Errorf("an unknown error ocurred") // want "`ocurred` is a misspelling of `occurred`"
}
