//args: -Epredeclared
//config_path: testdata/configs/predeclared.yml
package testdata

func hello() {
	var real int

	a := A{}
	println(a.true)

	copy := Clone(a)

	// suppress any "declared but not used" errors
	_ = real
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

func recover() {}
