//golangcitest:args -Egrouper
//golangcitest:config_path testdata/grouper.yml
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C" // want "should only use grouped 'import' declarations"

import "unsafe" // want "grouper\\(related information\\): found here"

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
