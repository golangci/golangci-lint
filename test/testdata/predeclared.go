//args: -Epredeclared
package testdata

func hello() {
	var real int // ERROR "variable real has same name as predeclared identifier"
	a := A{}
	copy := Clone(a) // ERROR "variable copy has same name as predeclared identifier"

	// suppress any "declared but not used" errors
	_ = real
	_ = a
	_ = copy
}

type A struct {
	true bool
	foo  int
}

func Clone(a A) A {
	return A{
		true: a.true,
		foo:  a.foo,
	}
}

func recover() {} // ERROR "function recover has same name as predeclared identifier"
