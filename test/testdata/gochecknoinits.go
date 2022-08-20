//golangcitest:args -Egochecknoinits
package testdata

import "fmt"

func init() { // want "don't use `init` function"
	fmt.Println()
}

func Init() {}
