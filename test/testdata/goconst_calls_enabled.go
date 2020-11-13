//args: -Egoconst
//config: linters-settings.goconst.ignore-calls=false
package testdata

import "fmt"

const FooBar = "foobar"

func Baz() {
	a := "foobar" // ERROR "string `foobar` has 3 occurrences as assignment statement, but such constant `FooBar` already exists"
	fmt.Print(a)
	b := "foobar"
	fmt.Print(b)
	c := "foobar"
	fmt.Print(c)
	fmt.Print("foobar") // ERROR "string `foobar` has 1 occurrences as call statement, but such constant `FooBar` already exists"
}
