//golangcitest:args -Eineffassign
package testdata

import "math"

func _() {
	x := math.MinInt8
	for {
		_ = x
		x = 0 // want "ineffectual assignment to x"
		x = 0
	}
}
