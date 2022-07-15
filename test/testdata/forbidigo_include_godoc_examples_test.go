//golangcitest:args -Eforbidigo
//golangcitest:config linters-settings.forbidigo.exclude-godoc-examples=false
package testdata

import "fmt"

func ExampleForbidigoNoGodoc() {
	fmt.Printf("too noisy!!!") // ERROR "use of `fmt.Printf` forbidden by pattern.*"
}
