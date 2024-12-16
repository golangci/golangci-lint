//golangcitest:args -Eexportloopref
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

func dummyFunction() {
	var array [4]*int
	var slice []*int
	var ref *int
	var str struct{ x *int }

	fmt.Println("loop expecting 10, 11, 12, 13")
	for i, p := range []int{10, 11, 12, 13} {
		printp(&p)
		slice = append(slice, &p) // want "exporting a pointer for the loop variable p"
		array[i] = &p             // want "exporting a pointer for the loop variable p"
		if i%2 == 0 {
			ref = &p   // want "exporting a pointer for the loop variable p"
			str.x = &p // want "exporting a pointer for the loop variable p"
		}
		var vStr struct{ x *int }
		var vArray [4]*int
		var v *int
		if i%2 == 0 {
			v = &p
			vArray[1] = &p
			vStr.x = &p
		}
		_ = v
	}

	fmt.Println(`slice expecting "10, 11, 12, 13" but "13, 13, 13, 13"`)
	for _, p := range slice {
		printp(p)
	}
	fmt.Println(`array expecting "10, 11, 12, 13" but "13, 13, 13, 13"`)
	for _, p := range array {
		printp(p)
	}
	fmt.Println(`captured value expecting "12" but "13"`)
	printp(ref)
}

func printp(p *int) {
	fmt.Println(*p)
}
