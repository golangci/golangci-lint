//golangcitest:args -Egoptrcmp
package testdata

import (
	"fmt"
)

type Foo struct {
	String  string
	Pointer *string
}

func CmpFoos(a, b *Foo) bool {
	if a != b { // want "pointer comparison: a != b"
		return false
	}

	if a == nil {
		return false
	}

	if nil == b {
		return false
	}

	if *a != *b {
		return false
	}

	if a.String != b.String {
		return false
	}

	if a.Pointer != b.Pointer { // want "pointer comparison: a.Pointer != b.Pointer"
		return false
	}

	if *a.Pointer != *b.Pointer {
		return false
	}

	fmt.Println(a == b) // want "pointer comparison: a == b"
	fmt.Println(a.String == b.String)
	fmt.Println(a.Pointer == b.Pointer) // want "pointer comparison: a.Pointer == b.Pointer"
	fmt.Println(*a.Pointer == *b.Pointer)

	return true
}
