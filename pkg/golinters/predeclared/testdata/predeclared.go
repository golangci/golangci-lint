//golangcitest:args -Epredeclared
package testdata

func hello() {
	var real int // want "variable real has same name as predeclared identifier"
	a := A{}
	copy := Clone(a) // want "variable copy has same name as predeclared identifier"

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

func recover() {} // want "function recover has same name as predeclared identifier"
