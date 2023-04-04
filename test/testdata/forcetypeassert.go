//golangcitest:args -Eforcetypeassert
package testdata

import "fmt"

func forcetypeassertInvalid() {
	var a interface{}
	_ = a.(int) // want "type assertion must be checked"

	var b interface{}
	bi := b.(int) // want "type assertion must be checked"
	fmt.Println(bi)
}

func forcetypeassertValid() {
	var a interface{}
	if ai, ok := a.(int); ok {
		fmt.Println(ai)
	}
}
