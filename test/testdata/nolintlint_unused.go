//args: -Enolintlint -Evarcheck
//config: linters-settings.nolintlint.allow-unused=false
package testdata

import "fmt"

func Foo() {
	fmt.Println("unused")          //nolint // ERROR "directive `//nolint .*` is unused"
	fmt.Println("unused,specific") //nolint:varcheck // ERROR "directive `//nolint:varcheck .*` is unused for linter varcheck"
	fmt.Println("not run")         //nolint:unparam // unparam is not run so this is ok
}
