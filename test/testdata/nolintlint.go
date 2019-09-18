//args: -Enolintlint
//config: linters-settings.nolintlint.explain=true
//config: linters-settings.nolintlint.specific=true
//config: linters-settings.nolintlint.machine=true
package testdata

import "fmt"

func Foo() {
	fmt.Println("not specific")         //nolint // ERROR "directive `.*` should mention specific linter such as `//nolint:my-linter`"
	fmt.Println("not machine readable") // nolint // ERROR "directive `.*`  should be written as `//nolint`"
	fmt.Println("bad syntax")           //nolint: deadcode // ERROR "directive `.*` should match `//nolint\[:<comma-separated-linters>\] \[// <explanation>\]`"
	fmt.Println("bad syntax")           //nolint:deadcode lll // ERROR "directive `.*` should match `//nolint\[:<comma-separated-linters>\] \[// <explanation>\]`"
	fmt.Println("extra spaces")         //  nolint:deadcode // because // ERROR "directive `.*` should not have more than one leading space"
}
