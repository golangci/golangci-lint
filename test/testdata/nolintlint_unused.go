//golangcitest:args -Enolintlint -Evarcheck
//golangcitest:config linters-settings.nolintlint.allow-unused=false
//golangcitest:expected_linter nolintlint
package testdata

import "fmt"

func Foo() {
	fmt.Println("unused")          //nolint:all // ERROR "directive `//nolint .*` is unused"
	fmt.Println("unused,specific") //nolint:varcheck // ERROR "directive `//nolint:varcheck .*` is unused for linter varcheck"
	fmt.Println("not run")         //nolint:unparam // unparam is not run so this is ok
}
