//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package minmax

func ifmin(a, b int) {
	x := a // A
	// B
	if a < b { // want "if statement can be modernized using max"
		// C
		x = b // D
		// E
	}
	print(x)
}

func ifmax(a, b int) {
	x := a
	if a > b { // want "if statement can be modernized using min"
		x = b
	}
	print(x)
}

func ifminvariant(a, b int) {
	x := a
	if x > b { // want "if statement can be modernized using min"
		x = b
	}
	print(x)
}

func ifmaxvariant(a, b int) {
	x := b
	if a < x { // want "if statement can be modernized using min"
		x = a
	}
	print(x)
}

func ifelsemin(a, b int) {
	var x int // A
	// B
	if a <= b { // want "if/else statement can be modernized using min"
		// C
		x = a // D
		// E
	} else {
		// F
		x = b // G
		// H
	}
	print(x)
}

func ifelsemax(a, b int) {
	// A
	var x int // B
	// C
	if a >= b { // want "if/else statement can be modernized using max"
		// D
		x = a // E
		// F
	} else {
		// G
		x = b
	}
	print(x)
}

func shadowed() int {
	hour, min := 3600, 60

	var time int
	if hour < min { // silent: the built-in min function is shadowed here
		time = hour
	} else {
		time = min
	}
	return time
}

func nopeIfStmtHasInitStmt() {
	x := 1
	if y := 2; y < x { // silent: IfStmt has an Init stmt
		x = y
	}
	print(x)
}

// Regression test for a bug: fix was "y := max(x, y)".
func oops() {
	x := 1
	y := 2
	if x > y { // want "if statement can be modernized using max"
		y = x
	}
	print(y)
}

// Regression test for a bug: += is not a simple assignment.
func nopeAssignHasIncrementOperator() {
	x := 1
	y := 0
	y += 2
	if x > y {
		y = x
	}
	print(y)
}

// Regression test for https://github.com/golang/go/issues/71721.
func nopeNotAMinimum(x, y int) int {
	// A value of -1 or 0 will use a default value (30).
	if x <= 0 {
		y = 30
	} else {
		y = x
	}
	return y
}

// Regression test for https://github.com/golang/go/issues/71847#issuecomment-2673491596
func nopeHasElseBlock(x int) int {
	y := x
	// Before, this was erroneously reduced to y = max(x, 0)
	if y < 0 {
		y = 0
	} else {
		y += 2
	}
	return y
}

func fix72727(a, b int) {
	o := a - 42
	// some important comment. DO NOT REMOVE.
	if o < b { // want "if statement can be modernized using max"
		o = b
	}
}

type myfloat float64

// The built-in min/max differ in their treatment of NaN,
// so reject floating-point numbers (#72829).
func nopeFloat(a, b myfloat) (res myfloat) {
	if a < b {
		res = a
	} else {
		res = b
	}
	return
}

// Regression test for golang/go#72928.
func underscoreAssign(a, b int) {
	if a > b {
		_ = a
	}
}

// Regression test for https://github.com/golang/go/issues/73576.
func nopeIfElseIf(a int) int {
	x := 0
	if a < 0 {
		x = 0
	} else if a > 100 {
		x = 100
	} else {
		x = a
	}
	return x
}
