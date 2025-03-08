//golangcitest:config_path testdata/gofumpt.yml
package testdata

import "fmt"

func GofumptNewLine() {
	fmt.Println( "foo" ) // want "File is not properly formatted"
}
