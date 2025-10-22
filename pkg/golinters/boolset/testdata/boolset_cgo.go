//golangcitest:args -Eboolset
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
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

var allNamesCgo = []string{"a", "a", "b", "c", "c", "d"}

func findDuplicatesCgo() []string {
	var duplicates []string
	uniqueNames := make(map[string]bool) // want "map\\[string\\]bool only stores \"true\" values; consider map\\[string\\]struct\\{\\}"
	for _, name := range allNamesCgo {
		if _, ok := uniqueNames[name]; ok {
			duplicates = append(duplicates, name)
			log.Println("Duplicate found: ", name)
		} else {
			uniqueNames[name] = true
		}
	}

	return duplicates
}
