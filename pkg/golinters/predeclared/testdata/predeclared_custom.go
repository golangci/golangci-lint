//golangcitest:args -Epredeclared
//golangcitest:config_path testdata/predeclared_custom.yml
package testdata

func hello() {
	var real int
	a := A{}
	copy := Clone(a) // want "copy: same name as predeclared identifier"

	// suppress any "declared but not used" errors
	_ = real
	_ = a
	_ = copy
}

type A struct {
	true bool // want "true: same name as predeclared identifier"
	foo  int
}

func Clone(a A) A {
	return A{
		true: a.true,
		foo:  a.foo,
	}
}

func recover() {}

type t1 byte

func (byte *t1) m1() {} // want "byte: same name as predeclared identifier"
func (t *t1) byte()  {} // want "byte: same name as predeclared identifier"
