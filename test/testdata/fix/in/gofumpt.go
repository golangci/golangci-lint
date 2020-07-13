//args: -Egofumpt
//config: linters-settings.gofumpt.extra-rules=true
package testdata

import "fmt"

func GofmtNotExtra(bar string, baz string) {
	fmt.Print(bar, baz)
}
