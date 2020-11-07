//args: -Eexhaustivestruct
package testdata

type Test struct {
	A string
	B int
	c bool // ignore private field
	D float64
}

var pass = Test{
	A: "a",
	B: 0,
	D: 1.0,
}

var fail = Test{ // ERROR "B is missing in Test"
	A: "a",
	D: 1.0,
}

var failMultiple = Test{ // ERROR "B, D are missing in Test"
	A: "a",
}
