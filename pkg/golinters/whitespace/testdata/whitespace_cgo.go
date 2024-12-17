//go:build ignore

// TODO(ldez) the linter doesn't support cgo.

//golangcitest:args -Ewhitespace
//golangcitest:config_path testdata/whitespace.yml
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

func oneLeadingNewline() { // want "unnecessary leading newline"

	fmt.Println("Hello world")
}

func oneNewlineAtBothEnds() { // want "unnecessary leading newline"

	fmt.Println("Hello world")

} // want "unnecessary trailing newline"

func noNewlineFunc() {
}

func oneNewlineFunc() {

}

func twoNewlinesFunc() {


}

func noNewlineWithCommentFunc() {
	// some comment
}

func oneTrailingNewlineWithCommentFunc() {
	// some comment

}

func oneLeadingNewlineWithCommentFunc() {

	// some comment
}

func twoLeadingNewlines() { // want "unnecessary leading newline"


	fmt.Println("Hello world")
}

func multiFuncFunc(a int,
	b int) { // want "multi-line statement should be followed by a newline"
	fmt.Println("Hello world")
}

func multiIfFunc() {
	if 1 == 1 &&
		2 == 2 { // want "multi-line statement should be followed by a newline"
		fmt.Println("Hello multi-line world")
	}

	if true {
		if true {
			if true {
				if 1 == 1 &&
					2 == 2 { // want "multi-line statement should be followed by a newline"
						fmt.Println("Hello nested multi-line world")
				}
			}
		}
	}
}

func notGoFmted() { // want "unnecessary leading newline"




         fmt.Println("Hello world")



} // want "unnecessary trailing newline"
