//args: -Eexhaustivestruct
package testdata

type Test struct {
	A string
	B int
}

var pass = Test{
	A: "a",
	B: 0,
}

var fail = Test{ // ERROR "B is missing in Test"
	A: "a",
}
