//golangcitest:args -Efilen
//golangcitest:config_path testdata/filen_min_lines.yml
package testdata // want "The number of lines in the file filen_min_lines.go less the allowed value! minLinesNum = 100, fileLines = 17"

import "fmt"

// foo0 test function
func foo0() string {
	fmt.Println("foo0 is a test")
	return "bar"
}

// foo1 test function
func foo1() string {
	fmt.Println("foo1 is a test")
	return "bar"
}
