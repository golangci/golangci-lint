//golangcitest:args -Egomodguard
//golangcitest:config_path testdata/gomodguard.yml
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

	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3" // want "import of package `gopkg.in/yaml.v3` is blocked because the module is in the blocked modules list. `github.com/kylelemons/go-gypsy` is a recommended module. This is an example of recommendations."
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

// Something just some struct
type Something struct{}

func aAllowedImport() { //nolint:unused
	mfile, _ := modfile.Parse("go.mod", []byte{}, nil)

	log.Println(mfile)
}

func aBlockedImport() { //nolint:unused
	data := []byte{}
	something := Something{}
	_ = yaml.Unmarshal(data, &something)

	log.Println(data)
}
