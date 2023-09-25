//golangcitest:args -Edupword
//golangcitest:config_path testdata/configs/dupword_ignore_the.yml
package testdata

import "fmt"

func duplicateWordInComments() {
	// this line include duplicated word the the
	fmt.Println("hello")
}

func duplicateWordInStr() {
	a := "this line include duplicate word and and"                   // want `Duplicate words \(and\) found`
	b := "print the\n the line, print the the \n\t the line. and and" // want `Duplicate words \(and\) found`
	fmt.Println(a, b)
}
