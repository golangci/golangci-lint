//golangcitest:args -Edupword
package testdata

import "fmt"

func duplicateWordInComments() {
	// this line include duplicated word the the // want `Duplicate words \(the\) found`
	fmt.Println("hello")
}

func duplicateWordInStr() {
	a := "this line include duplicate word and and"                   // want `Duplicate words \(and\) found`
	b := "print the\n the line, print the the \n\t the line. and and" // want `Duplicate words \(and,the\) found`
	fmt.Println(a, b)
}
