//golangcitest:args -Emakezero
//golangcitest:config_path testdata/configs/makezero_always.yml
package testdata

import "math"

func MakezeroAlways() []int {
	x := make([]int, math.MaxInt8) // want "slice `x` does not have non-zero initial length"
	return x
}

func MakezeroAlwaysNolint() []int {
	x := make([]int, math.MaxInt8) //nolint:makezero // ok that this is not initialized
	return x
}
