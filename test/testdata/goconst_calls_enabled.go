//golangcitest:args -Egoconst
//golangcitest:config_path testdata/configs/goconst_calls_enabled.yml
package testdata

import "fmt"

const FooBar = "foobar"

func Baz() {
	a := "foobar" // want "string `foobar` has 4 occurrences, but such constant `FooBar` already exists"
	fmt.Print(a)
	b := "foobar"
	fmt.Print(b)
	c := "foobar"
	fmt.Print(c)
	fmt.Print("foobar")
}
