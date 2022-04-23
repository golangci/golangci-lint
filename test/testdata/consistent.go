//args: -Econsistent
package testdata

import ( // ERROR "remove parens around single import declaration"
	"strings"
)

func consistent() {
	_ = func(a, b int) {} // ERROR "declare the type of function arguments explicitly"

	_ = func() (a, b int) { return 1, 2 } // ERROR "declare the type of function return values explicitly"

	_ = new(strings.Builder) // ERROR "use zero-value literal instead of calling new"

	_ = make([]int, 0) // ERROR "use slice literal instead of calling make"

	_ = 0xABCDE // ERROR "use lowercase digits in hex literal"

	x := 5
	_ = 1 < x && x < 10 // ERROR "write common term in range expression on the left"

	_ = 1 & ^2 // ERROR "use AND-NOT operator instead of AND operator with complement expression"

	_ = .5 // ERROR "add zero before decimal point in floating-point literal"

	_ = len([]int{}) > 0 // ERROR `check if len is \(not\) 0 instead`

	switch {
	case 1 < 2 || 3 < 4: // ERROR "separate cases with comma instead of using logical OR"
	}

	switch {
	default: // ERROR "move switch default clause to the end"
	case 1 < 2:
	}

test_loop: // ERROR `change label to match regular expression: .*`
	for {
		break test_loop
	}
}
