package testdata

import (
	"fmt"
)

const half = 0.5

// Deprecated: use floatRounder instead
func deprecatedFloatRounder(f float64) int64 {
	return int64(f)
}

func floatRounder(f float64) int64 {
	sgn := int64(1)
	if f < 0 {
		sgn = -1
	}
	out := int64(f)
	if f-float64(out) > half {
		out++
	}
	return sgn * out
}

func StaticCheckExitSuccessOnInfo() {
	fmt.Println(deprecatedFloatRounder(1.0)) // want "SA1019: deprecatedFloatRounder is deprecated"
	fmt.Println(floatRounder(1.0))
}
