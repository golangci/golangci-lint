//golangcitest:args -Egofumpt
//golangcitest:config linters-settings.gofumpt.extra-rules=true
package p

import "fmt"

func GofmtNotExtra(bar string, baz string) {
	fmt.Print(bar, baz)
}
