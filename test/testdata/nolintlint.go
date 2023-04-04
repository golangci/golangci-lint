//golangcitest:args -Enolintlint -Emisspell
//golangcitest:expected_linter nolintlint
//golangcitest:config_path testdata/configs/nolintlint.yml
package testdata

import "fmt"

func Foo() {
	fmt.Println("not specific")         //nolint // want "directive `.*` should mention specific linter such as `//nolint:my-linter`"
	fmt.Println("not machine readable") // nolint // want "directive `.*`  should be written as `//nolint`"
	fmt.Println("extra spaces")         //  nolint:unused // because // want "directive `.*` should not have more than one leading space"

	// test expanded range
	//nolint:misspell // deliberate misspelling to trigger nolintlint
	func() {
		mispell := true
		fmt.Println(mispell)
	}()
}
