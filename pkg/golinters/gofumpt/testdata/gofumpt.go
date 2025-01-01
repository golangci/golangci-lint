//golangcitest:args -Egofumpt
package testdata

import "fmt"

func GofumptNewLine() {
	fmt.Println( "foo" ) // want "File is not properly formatted"
}