//golangcitest:args -Eexhaustivestruct --internal-cmd-test
//golangcitest:config_path testdata/configs/exhaustivestruct.yml
package testdata

import "time"

type ExhaustiveStructCustom struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustiveStructCustom() {
	// pass
	_ = ExhaustiveStructCustom{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = ExhaustiveStructCustom{ // want "B is missing in ExhaustiveStructCustom"
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = ExhaustiveStructCustom{ // want "B, D are missing in ExhaustiveStructCustom"
		A: "a",
		c: false,
		E: time.Now(),
	}

	// failPrivate
	_ = ExhaustiveStructCustom{ // want "c is missing in ExhaustiveStructCustom"
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

}

type ExhaustiveStructCustom1 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustiveStructCustom1() {
	_ = ExhaustiveStructCustom1{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	_ = ExhaustiveStructCustom1{
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	_ = ExhaustiveStructCustom1{
		A: "a",
		c: false,
		E: time.Now(),
	}

	_ = ExhaustiveStructCustom1{
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}
}

type ExhaustiveStructCustom2 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustiveStructCustom2() {
	// pass
	_ = ExhaustiveStructCustom2{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = ExhaustiveStructCustom2{ // want "B is missing in ExhaustiveStructCustom2"
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = ExhaustiveStructCustom2{ // want "B, D are missing in ExhaustiveStructCustom2"
		A: "a",
		c: false,
		E: time.Now(),
	}

	// failPrivate
	_ = ExhaustiveStructCustom2{ // want "c is missing in ExhaustiveStructCustom2"
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

}
