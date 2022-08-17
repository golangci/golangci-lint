//golangcitest:args -Eexhaustruct
//golangcitest:config_path testdata/configs/exhaustruct.yml
package testdata

import "time"

type ExhaustructCustom struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustructCustom() {
	// pass
	_ = ExhaustructCustom{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = ExhaustructCustom{ // want "B is missing in ExhaustructCustom"
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = ExhaustructCustom{ // want "B, D are missing in ExhaustructCustom"
		A: "a",
		c: false,
		E: time.Now(),
	}

	//  failPrivate
	_ = ExhaustructCustom{ // want "c is missing in ExhaustructCustom"
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

}

type ExhaustructCustom1 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustructCustom1() {
	// pass
	_ = ExhaustructCustom{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = ExhaustructCustom1{
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = ExhaustructCustom1{
		A: "a",
		c: false,
		E: time.Now(),
	}

	// failPrivate
	_ = ExhaustructCustom1{
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}

}

type ExhaustructCustom2 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

func exhaustructCustom2() {
	// pass
	_ = ExhaustructCustom2{
		A: "a",
		B: 0,
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// fail
	_ = ExhaustructCustom2{
		A: "a",
		c: false,
		D: 1.0,
		E: time.Now(),
	}

	// failMultiple
	_ = ExhaustructCustom2{
		A: "a",
		c: false,
		E: time.Now(),
	}

	// failPrivate
	_ = ExhaustructCustom2{
		A: "a",
		B: 0,
		D: 1.0,
		E: time.Now(),
	}
}
