//golangcitest:args -Egci
//golangcitest:config_path testdata/gci.yml
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
	"golang.org/x/tools/go/analysis" // want "File is not properly formatted"
	"github.com/golangci/golangci-lint/pkg/config"
	"fmt"
	"errors"
	gcicfg "github.com/daixiang0/gci/pkg/config"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func GoimportsLocalTest() {
	fmt.Print(errors.New("x"))
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
