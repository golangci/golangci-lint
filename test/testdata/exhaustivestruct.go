//args: -Eexhaustivestruct
package testdata

type Test struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
}

var pass = Test{
	A: "a",
	B: 0,
	c: false,
	D: 1.0,
}

var failPrivate = Test{ // ERROR "c is missing in Test"
	A: "a",
	B: 0,
	D: 1.0,
}

var fail = Test{ // ERROR "B is missing in Test"
	A: "a",
	c: false,
	D: 1.0,
}

var failMultiple = Test{ // ERROR "B, D are missing in Test"
	A: "a",
	c: false,
}
