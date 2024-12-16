//golangcitest:args -Emirror
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
	"strings"
	"unicode/utf8"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	_ = utf8.RuneCount([]byte("foobar"))                                                                            // want `avoid allocations with utf8\.RuneCountInString`
	_ = strings.Compare(string([]byte{'f', 'o', 'o', 'b', 'a', 'r'}), string([]byte{'f', 'o', 'o', 'b', 'a', 'r'})) // want `avoid allocations with bytes\.Compare`
}
