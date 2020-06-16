//args: -Egofumpt
package testdata

import "fmt"

func GofumptNewLine() {

	fmt.Println("foo")
}

// ERROR "File is not `gofumpt`-ed"
