//golangcitest:args -Eprealloc
package testdata

func Prealloc(source []int) []int {
	var dest []int // want `Consider preallocating dest with capacity len\(source\)`
	for _, v := range source {
		dest = append(dest, v)
	}

	return dest
}
