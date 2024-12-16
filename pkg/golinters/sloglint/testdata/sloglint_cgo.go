//go:build go1.21

//golangcitest:args -Esloglint
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
	"log/slog"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	slog.Info("msg", "foo", 1, "bar", 2)
	slog.Info("msg", slog.Int("foo", 1), slog.Int("bar", 2))

	slog.Info("msg", "foo", 1, slog.Int("bar", 2)) // want `key-value pairs and attributes should not be mixed`
}
