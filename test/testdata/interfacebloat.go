//golangcitest:args -Einterfacebloat
package testdata

import "time"

type InterfaceBloatExample01 interface { // want "the interface has more than 10 methods: 11"
	a01() time.Duration
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
}

func InterfaceBloatExample02() {
	var _ interface { // want "the interface has more than 10 methods: 11"
		a01() time.Duration
		a02()
		a03()
		a04()
		a05()
		a06()
		a07()
		a08()
		a09()
		a10()
		a11()
	}
}

func InterfaceBloatExample03() interface { // want "the interface has more than 10 methods: 11"
	a01() time.Duration
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
} {
	return nil
}

type InterfaceBloatExample04 struct {
	Foo interface { // want "the interface has more than 10 methods: 11"
		a01() time.Duration
		a02()
		a03()
		a04()
		a05()
		a06()
		a07()
		a08()
		a09()
		a10()
		a11()
	}
}

type InterfaceBloatSmall01 interface {
	a01() time.Duration
	a02()
	a03()
	a04()
	a05()
}

type InterfaceBloatSmall02 interface {
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
}

type InterfaceBloatExample05 interface {
	InterfaceBloatSmall01
	InterfaceBloatSmall02
}

type InterfaceBloatExample06 interface {
	interface { // want "the interface has more than 10 methods: 11"
		a01() time.Duration
		a02()
		a03()
		a04()
		a05()
		a06()
		a07()
		a08()
		a09()
		a10()
		a11()
	}
}

type InterfaceBloatTypeGeneric interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | uint |
		~int8 | ~int16 | ~int32 | ~int64 | int |
		~float32 | ~float64 |
		~string
}

func InterfaceBloatExampleNoProblem() interface {
	a01() time.Duration
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
} {
	return nil
}
