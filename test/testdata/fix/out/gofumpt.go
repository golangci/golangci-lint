//golangcitest:args -Egofumpt
//golangcitest:config_path testdata/configs/gofumpt-fix.yml
package p

import "fmt"

func GofmtNotExtra(bar, baz string) {
	fmt.Print(bar, baz)
}
