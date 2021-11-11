//args: -Evarnamelen
package testdata

import "math"

func varnamelen() {
	x := math.MinInt8 // ERROR "variable name 'x' is too short for the scope of its usage"
	x++
	x++
	x++
	x++
	x++
	x++
	x++
	x++
	x++
}
