//golangcitest:args -Edepguard
//golangcitest:config_path testdata/depguard.yml
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
	"compress/gzip" // want "import 'compress/gzip' is not allowed from list 'main': nope"
	"log"           // want "import 'log' is not allowed from list 'main': don't use log"
	"unsafe"

	"golang.org/x/tools/go/analysis" // want "import 'golang.org/x/tools/go/analysis' is not allowed from list 'main': example import with dot"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
	_ = analysis.Analyzer{}
}
