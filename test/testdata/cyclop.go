// args: -Ecyclop
package testdata

import "math"

func cyclopComplexFunc() { // ERROR "calculated cyclomatic complexity for function cyclopComplexFunc is 11, max is 10"
	i := math.MaxInt8
	if i > 2 {
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
	} else {
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
	}

	if i > 2 {
	}
}
