//golangcitest:args -Eenonly
package testdata

import (
	"fmt"
)

var v = "فارسی" // want "contains non-english characters: \"فارسی\""

const v2 = "فارسی" // want "contains non-english characters: \"فارسی\""

func example() {
	var test1 = "فارسی" // want "contains non-english characters: \"فارسی\""
	_ = test1

	test2 := "فارسی" // want "contains non-english characters: \"فارسی\""
	_ = test2

	var test3 string
	test3 = "فارسی" // want "contains non-english characters: \"فارسی\""
	_ = test3

	example2("فارسی") // want "contains non-english characters: \"فارسی\""
}

func example2(s string) {
	fmt.Println(s)
}
