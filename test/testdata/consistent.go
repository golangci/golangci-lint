//golangcitest:args -Econsistent
package testdata

import ( // want "remove parens around single import declaration"
	"strings"
)

func consistent() {
	_ = func(a, b int) {} // want "declare the type of function arguments explicitly"

	_ = func() (a, b int) { return 1, 2 } // want "declare the type of function return values explicitly"

	_ = new(strings.Builder) // want "use zero-value literal instead of calling new"

	_ = make([]int, 0) // want "use slice literal instead of calling make"

	_ = 0xABCDE // want "use lowercase digits in hex literal"

	x := 5
	_ = 1 < x && x < 10 // want "write common term in range expression on the left"

	_ = 1 & ^2 // want "use AND-NOT operator instead of AND operator with complement expression"

	_ = .5 // want "add zero before decimal point in floating-point literal"

	_ = len([]int{}) > 0 // want `check if len is \(not\) 0 instead`

	switch {
	case 1 < 2 || 3 < 4: // want "separate cases with comma instead of using logical OR"
	}

	switch {
	default: // want "move switch default clause to the end"
	case 1 < 2:
	}

	type empty interface{} // want "use any instead of interface{}"

test_loop: // want `change label to match regular expression: .*`
	for {
		break test_loop
	}
}
