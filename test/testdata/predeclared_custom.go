//golangcitest:args -Epredeclared
//golangcitest:config_path testdata/configs/predeclared.yml
package testdata

func hello() {
	var real int
	a := A{}
	copy := Clone(a) // ERROR "variable copy has same name as predeclared identifier"

	// suppress any "declared but not used" errors
	_ = real
	_ = a
	_ = copy
}

type A struct {
	true bool // ERROR "field true has same name as predeclared identifier"
	foo  int
}

func Clone(a A) A {
	return A{
		true: a.true,
		foo:  a.foo,
	}
}

func recover() {}
