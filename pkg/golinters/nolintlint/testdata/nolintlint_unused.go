//golangcitest:args -Enolintlint -Evarcheck
//golangcitest:expected_linter nolintlint
//golangcitest:config_path nolintlint_unused.yml
package testdata

import "fmt"

func Foo() {
	fmt.Println("unused")          //nolint:all // want "directive `//nolint .*` is unused"
	fmt.Println("unused,specific") //nolint:varcheck // want "directive `//nolint:varcheck .*` is unused for linter varcheck"
	fmt.Println("not run")         //nolint:unparam // unparam is not run so this is ok
}
