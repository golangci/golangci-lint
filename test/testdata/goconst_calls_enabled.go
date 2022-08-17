//golangcitest:args -Egoconst
//golangcitest:config linters-settings.goconst.ignore-calls=false
package testdata

import "fmt"

const FooBar = "foobar"

func Baz() {
	a := "foobar" // ERROR "string `foobar` has 4 occurrences, but such constant `FooBar` already exists"
	fmt.Print(a)
	b := "foobar"
	fmt.Print(b)
	c := "foobar"
	fmt.Print(c)
	fmt.Print("foobar")
}
