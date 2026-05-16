//golangcitest:args -Egoconst
//golangcitest:config_path testdata/goconst_exclude_types.yml
package testdata

import "fmt"

const FooBar = "foobar"

func _() {
	a := "foobar" // want "string `foobar` has 3 occurrences, but such constant `FooBar` already exists"
	fmt.Print(a)

	_ = []string{"foobar"}

	_ = map[string]string{"foobar": "value"}

	b := "foobar"
	fmt.Print(b)

	c := "foobar"
	fmt.Print(c)

	fmt.Print("foobar")

	fmt.Print("fdsqfdsqfezrazt")
}
