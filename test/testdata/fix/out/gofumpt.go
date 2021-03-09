// args: -Egofumpt
// config: linters-settings.gofumpt.extra-rules=true
package p

import "fmt"

func GofmtNotExtra(bar, baz string) {
	fmt.Print(bar, baz)
}
