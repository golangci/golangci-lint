//golangcitest:args -Eexhaustivestruct --internal-cmd-test
package testdata

import "time"

type ExhaustiveStruct struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustiveStruct() {
	// pass
	_ = ExhaustiveStruct{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failPrivate
	_ = ExhaustiveStruct{ // ERROR "c is missing in ExhaustiveStruct"
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = ExhaustiveStruct{ // ERROR "B is missing in ExhaustiveStruct"
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = ExhaustiveStruct{ // ERROR "B, D are missing in ExhaustiveStruct"
		A: "a",
		c: false,
		E: time.Now(),
	}
}
