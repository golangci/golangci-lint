//golangcitest:args -Erevive
//golangcitest:config_path testdata/revive_add_constant.yml
package testdata

import "fmt"

func testReviveAddConstant() {
	fmt.Println("dup")
	fmt.Println("dup")
	fmt.Println("dup") // want `add-constant: string literal "dup" appears, at least, 3 times, create a named constant for it`
	fmt.Println(42)    // want `add-constant: avoid magic numbers like '42', create a named constant for it`
	fmt.Println(1.0)
	fmt.Println("")
}
