//golangcitest:args -Econstructorcheck
package testdata

import (
	"bytes"
	"fmt"
)

var buf = bytes.Buffer{} // standard library is excluded from analysis

// T is a type whose zero values are supposedly invalid
// so a constructor NewT was created.
type T struct { // want T:`{NewT \d* \d*}`
	x int
	s string
	m map[int]int
}

var (
	tNil       *T       // want `nil value of type p.T may be unsafe, use constructor NewT instead`
	tZero      = T{}    // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	tZeroPtr   = &T{}   // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	tNew       = new(T) // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	tComposite = T{     // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	tCompositePtr = &T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
		s: "abc",
	}
	tColl    = []T{T{x: 1}}   // want `use constructor NewT for type p.T instead of a composite literal`
	tPtrColl = []*T{&T{x: 1}} // want `use constructor NewT for type p.T instead of a composite literal`

)

// NewT is a valid constructor for type T. Here we check if it's called
// instead of constructing values of type T manually
func NewT() *T {
	return &T{
		m: make(map[int]int),
	}
}

type structWithTField struct {
	i int
	t T
}

var structWithT = structWithTField{
	i: 1,
	t: T{x: 1}, // want `use constructor NewT for type p.T instead of a composite literal`
}

type structWithTPtrField struct {
	i int
	t *T
}

var structWithTPtr = structWithTPtrField{
	i: 1,
	t: &T{x: 1}, // want `use constructor NewT for type p.T instead of a composite literal`
}

func fnWithT() {
	x := T{}     // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	x2 := &T{}   // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	x3 := new(T) // want `zero value of type p.T may be unsafe, use constructor NewT instead`
	fmt.Println(x, x2, x3)
}

func retT() T {
	return T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
	}
}

func retTPtr() *T {
	return &T{ // want `use constructor NewT for type p.T instead of a composite literal`
		x: 1,
	}
}

func retTNilPtr() *T {
	var t *T // want `nil value of type p.T may be unsafe, use constructor NewT instead`
	return t
}

type T2 struct { // want T2:`{NewT2 \d* \d*}`
	x int
}

func NewT2() *T2 {
	// new(T) inside T's constructor is permitted
	return new(T2)
}
