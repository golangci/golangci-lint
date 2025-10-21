//golangcitest:args -Edupword
//golangcitest:config_path testdata/dupword_comments_only.yml
package testdata

import "fmt"

func _() {
	// want +2 `Duplicate words \(and\) found`
	// want +2 `Duplicate words \(and,the\) found`
	// this line include duplicate word and and
	// print the\n the line, print the the \n\t the line. and and

	a := "this line include duplicate word and and"
	b := "print the\n the line, print the the \n\t the line. and and"
	fmt.Println(a, b)
}
