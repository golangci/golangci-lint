//args: -Emakezero
package testdata

func Makezero() []int {
	x := make([]int, 5)
	return append(x, 1)
}
