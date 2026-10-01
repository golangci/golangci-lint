//golangcitest:args -Epredeclared
package testdata

func hello() {
	var real int // want "real: shadows predeclared identifier"
	a := A{}
	copy := Clone(a) // want "copy: shadows predeclared identifier"

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

func recover() {} // want "recover: shadows predeclared identifier"

type t1 byte

func (byte *t1) m1() {} // want "byte: shadows predeclared identifier$"
func (t *t1) byte()  {}
