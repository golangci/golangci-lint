//golangcitest:args -Edupword
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func duplicateWordInComments() {
	// this line include duplicated word the the
	fmt.Println("hello")
}

func duplicateWordInStr() {
	a := "this line include duplicate word and and"
	b := "print the\n the line, print the the \n\t the line. and and"
	fmt.Println(a, b)
}
