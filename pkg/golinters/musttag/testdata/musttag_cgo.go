//golangcitest:args -Emusttag
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
	"encoding/asn1"
	"encoding/json"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

// builtin functions:
func _() {
	var user struct {
		Name  string
		Email string `json:"email"`
	}
	json.Marshal(user) // want "the given struct should be annotated with the `json` tag"
}

// custom functions from config:
func _() {
	var user struct {
		Name  string
		Email string `asn1:"email"`
	}
	asn1.Marshal(user)
}
