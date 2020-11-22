//args: -Emakezero
//config: linters-settings.makezero.always=true
package testdata

func MakezeroAlways() []int {
	x := make([]int, 5)
	return x // ERROR "slice `x` does not have non-zero initial length"
}


func MakezeroAlwaysNolint() []int {
	x := make([]int, 5)
	return x //nolint:makezero // ok that this is not initialized
}
