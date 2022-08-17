//golangcitest:args -Egofumpt
package testdata

import "fmt"

func GofumptNewLine() {
	fmt.Println( "foo" ) // want "File is not `gofumpt`-ed"
}
