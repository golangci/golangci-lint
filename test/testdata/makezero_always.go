//args: -Emakezero
//config: linters-settings.makezero.always=true
package testdata

func MakezeroAlways() []int {
	x := make([]int, 5) // ERROR "slice `x` does not have non-zero initial length"
	return x
}

func MakezeroAlwaysNolint() []int {
	x := make([]int, 5) //nolint:makezero // ok that this is not initialized
	return x
}
