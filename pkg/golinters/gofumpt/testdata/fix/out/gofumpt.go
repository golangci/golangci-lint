//golangcitest:config_path testdata/gofumpt-fix.yml
//golangcitest:expected_exitcode 0
package p

import "fmt"

func GofmtNotExtra(bar, baz string) {
	fmt.Print(bar, baz)
}
