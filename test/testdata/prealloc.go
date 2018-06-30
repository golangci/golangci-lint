// args: -Eprealloc
package testdata

func Prealloc(source []int) []int {
	var dest []int // ERROR "Consider preallocating `dest`"
	for _, v := range source {
		dest = append(dest, v)
	}

	return dest
}
