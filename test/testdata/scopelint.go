//args: -Escopelint
package testdata

import "fmt"

func ScopelintTest() {
	values := []string{"a", "b", "c"}
	var funcs []func()
	for _, val := range values {
		funcs = append(funcs, func() {
			fmt.Println(val) // ERROR "Using the variable on range scope `val` in function literal"
		})
	}
	for _, f := range funcs {
		f()
	}
}
