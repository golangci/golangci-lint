//golangcitest:args -Etypeswitch
package testdata

type I interface{ i() }

type A struct{}

func (A) i() {}

type B struct{}

func (B) i() {}

type C struct{}

func (C) i() {}

func g() {
	var i I = A{}

	switch i.(type) { // want "type C does not appear in any cases"
	case A:
	case B:
	}

	switch i.(type) { // want "type B does not appear in any cases"
	case A:
		println("a")
	case C:
		println("c")
	}

	switch v := i.(type) { // want "type B does not appear in any cases"
	case A:
		println(v)
	case C:
		println(v)
	}

	switch i.(type) { // OK
	case A, B, C:
		println("a, b, c")
	}

	switch i.(type) { // OK
	case A:
	default:
	}

	switch i.(type) { // OK
	case A:
	default:
	}
}
