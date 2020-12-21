//args: -Enolintlint -Evarcheck -Edeadcode
//config_path: testdata/configs/nolintlint_unused.yml
package testdata

import "fmt"

func Foo() {
	fmt.Println("unused")          // nolint // ERROR "directive `//nolint .*` is unused"
	fmt.Println("unused,specific") // nolint:varcheck // ERROR "directive `//nolint:varcheck .*` is unused for linter varcheck"
	fmt.Println("not run")         // nolint:unparam // unparam is not run so this is ok

	fmt.Println("unused but ignored by line") // nolint:varcheck,nolintlint // varcheck is unused but nolintlint is marked as nolint itself

	fmt.Println("unused but ignored by config") // nolint:deadcode // this linter is ignored for unused checking in config
}
