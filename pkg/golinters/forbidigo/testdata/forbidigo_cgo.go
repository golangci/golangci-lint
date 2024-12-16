//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/forbidigo.yml
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
	fmt2 "fmt"
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func Forbidigo() {
	fmt.Printf("too noisy!!!")  // want "use of `fmt\\.Printf` forbidden by pattern `fmt\\\\.Print\\.\\*`"
	fmt2.Printf("too noisy!!!") // Not detected because analyze-types is false by default for backward compatibility.
	time.Sleep(time.Nanosecond) // want "no sleeping!"
}
