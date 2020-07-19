//args: -Egofumpt
//config: linters-settings.gofumpt.extra-rules=true
package testdata

import "fmt"

func GofmtNotExtra(bar string, baz string) { // ERROR "File is not `gofumpt`-ed with `-extra`"
	fmt.Print("foo")
}
