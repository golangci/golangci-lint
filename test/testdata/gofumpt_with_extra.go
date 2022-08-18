//golangcitest:args -Egofumpt
//golangcitest:config_path testdata/configs/gofumpt_with_extra.yml
package testdata

import "fmt"

func GofmtNotExtra(bar string, baz string) { // ERROR "File is not `gofumpt`-ed with `-extra`"
	fmt.Print("foo")
}
