//golangcitest:args -Edupword
//golangcitest:config_path testdata/dupword_comments_only.yml
package testdata

import "fmt"

func duplicateWordInComments() {
	// this line include duplicated word the the // want `Duplicate words \(the\) found`
	fmt.Println("hello")
}

func duplicateWordInStr() {
	// When comments-only is enabled, duplicates in strings should NOT be flagged
	a := "this line include duplicate word and and"
	b := "print the\n the line, print the the \n\t the line. and and"
	fmt.Println(a, b)
}

// Another comment with with duplicate words // want `Duplicate words \(with\) found`
func anotherFunc() {
	// Duplicate in in comment should be caught // want `Duplicate words \(in\) found`
	s := "but but duplicate in string should be ignored"
	fmt.Println(s)
}
