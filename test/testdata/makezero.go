//args: -Emakezero
package testdata

func Makezero() []int {
	x := make([]int, 5)
	return append(x, 1) // ERROR "append to slice `x` with non-zero initialized length"
}

func MakezeroNolint() []int {
	x := make([]int, 5)
	return append(x, 1) //nolint:makezero // ok that we're appending to an uninitialized slice
}
