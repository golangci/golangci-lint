//args: -Emakezero
//config: linters-settings.makezero.always=true
package testdata

func MakezeroAlways() []int {
	return make([]int, 5)
}
