//args: -Egochecknoinits
package testdata

import "fmt"

func init() { // ERROR "don't use `init` function"
	fmt.Println()
}

func Init() {}
