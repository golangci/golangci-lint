//golangcitest:args -Eimportas
//golangcitest:config_path testdata/importas.yml
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
	wrong_alias "fmt" // want `import "fmt" imported as "wrong_alias" but must be "fff" according to config`
	"os"
	wrong_alias_again "os" // want `import "os" imported as "wrong_alias_again" but must be "std_os" according to config`
	"unsafe"

	wrong "golang.org/x/tools/go/analysis" // want `import "golang.org/x/tools/go/analysis" imported as "wrong" but must be "ananas" according to config`
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
	_ = wrong.Analyzer{}
}
