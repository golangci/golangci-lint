//golangcitest:args -Enolintlint -Emisspell
//golangcitest:expected_linter nolintlint
//golangcitest:config_path nolintlint_unused.yml
package testdata

import "fmt"

func Foo() {
	fmt.Println("unused")          //nolint:all // want "directive `//nolint .*` is unused"
	fmt.Println("unused,specific") //nolint:misspell // want "directive `//nolint:misspell .*` is unused for linter misspell"
	fmt.Println("not run")         //nolint:unparam // unparam is not run so this is ok
}
