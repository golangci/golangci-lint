//args: -Eexhaustivestruct
//config: linters-settings.exhaustivestruct.struct-patterns=*.Test1,*.Test3
package testdata

import "time"

type Test1 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

var passTest1 = Test1{
	A: "a",
	B: 0,
	c: false,
	D: 1.0,
	E: time.Now(),
}

var failTest1 = Test1{ // ERROR "B is missing in Test"
	A: "a",
	c: false,
	D: 1.0,
	E: time.Now(),
}

var failMultipleTest1 = Test1{ // ERROR "B, D are missing in Test"
	A: "a",
	c: false,
	E: time.Now(),
}

var failPrivateTest1 = Test1{ // ERROR "c is missing in Test"
	A: "a",
	B: 0,
	D: 1.0,
	E: time.Now(),
}

type Test2 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

var passTest2 = Test1{
	A: "a",
	B: 0,
	c: false,
	D: 1.0,
	E: time.Now(),
}

var failTest2 = Test2{
	A: "a",
	c: false,
	D: 1.0,
	E: time.Now(),
}

var failMultipleTest2 = Test2{
	A: "a",
	c: false,
	E: time.Now(),
}

var failPrivateTest2 = Test2{
	A: "a",
	B: 0,
	D: 1.0,
	E: time.Now(),
}

type Test3 struct {
	A string
	B int
	c bool // private field inside the same package are not ignored
	D float64
	E time.Time
}

var passTest3 = Test3{
	A: "a",
	B: 0,
	c: false,
	D: 1.0,
	E: time.Now(),
}

var failTest3 = Test3{ // ERROR "B is missing in Test"
	A: "a",
	c: false,
	D: 1.0,
	E: time.Now(),
}

var failMultipleTest3 = Test3{ // ERROR "B, D are missing in Test"
	A: "a",
	c: false,
	E: time.Now(),
}

var failPrivateTest3 = Test3{ // ERROR "c is missing in Test"
	A: "a",
	B: 0,
	D: 1.0,
	E: time.Now(),
}
