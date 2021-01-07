// args: -Egofumpt
// config: linters-settings.gofumpt.extra-rules=true
package testdata

import "fmt"

func GofmtNotExtra(bar, baz string) {
	fmt.Print(bar, baz)
}
