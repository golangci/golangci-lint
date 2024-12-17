//golangcitest:args -Enolintlint -Emisspell
//golangcitest:expected_linter nolintlint
//golangcitest:config_path nolintlint.yml
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

func Foo() {
	fmt.Println("not specific")         //nolint // want "directive `.*` should mention specific linter such as `//nolint:my-linter`"
	fmt.Println("not machine readable") // nolint // want "directive `.*`  should be written as `//nolint`"
	fmt.Println("extra spaces")         //  nolint:unused // because // want "directive `.*` should not have more than one leading space"

	// test expanded range
	//nolint:misspell // deliberate misspelling to trigger nolintlint
	func() {
		mispell := true
		fmt.Println(mispell)
	}()
}
