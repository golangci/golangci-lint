// args: -Eexhaustruct
package testdata

import "time"

type Exhaustruct struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustruct() {
	// pass
	_ = Exhaustruct{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failPrivate
	_ = Exhaustruct{ // ERROR "c is missing in Exhaustruct"
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = Exhaustruct{ // ERROR "B is missing in Exhaustruct"
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = Exhaustruct{ // ERROR "B, D are missing in Exhaustruct"
		A: "a",
		c: false,
		E: time.Now(),
	}

}
