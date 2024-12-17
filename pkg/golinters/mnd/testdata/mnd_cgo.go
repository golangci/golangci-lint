//golangcitest:args -Emnd
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
	"net/http"
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	c := &http.Client{
		Timeout: 2 * time.Second, // want "Magic number: 2, in <assign> detected"
	}

	res, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 { // want "Magic number: 200, in <condition> detected"
		log.Println("Something went wrong")
	}
}
