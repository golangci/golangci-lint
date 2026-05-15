//golangcitest:config_path testdata/gofumpt_with_extra.yml
package testdata

import "fmt"

func GroupedType(bar string, baz string) { // want "File is not properly formatted"
	fmt.Print("foo")
}

func RakedReturns() (err error) {
	return
}
