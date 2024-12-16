//golangcitest:args -Edurationcheck
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
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type durationCheckData struct {
	i int
	d time.Duration
}

func durationcheckCase01() {
	dcd := durationCheckData{i: 10}
	_ = time.Duration(dcd.i) * time.Second
}

func durationcheckCase02() {
	dcd := durationCheckData{d: 10 * time.Second}
	_ = dcd.d * time.Second // want "Multiplication of durations: `dcd.d \\* time.Second`"
}
