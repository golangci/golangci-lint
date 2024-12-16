//golangcitest:args -Edupword
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func duplicateWordInComments() {
	// this line include duplicated word the
	fmt.Println("hello")
}

func duplicateWordInStr() {
	a := "this line include duplicate word and"
	b := "print the\n line, print the line. and"
	fmt.Println(a, b)
}
