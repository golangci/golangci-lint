//golangcitest:args -Egofumpt
//golangcitest:config_path testdata/configs/gofumpt-fix.yml
//golangcitest:expected_exitcode 0
package p

import "fmt"

func GofmtNotExtra(bar string, baz string) {
	fmt.Print(bar, baz)
}
