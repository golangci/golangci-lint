// Package unique provides primitives for sorting slices removing repeated elements.
package unique

import (
	"reflect"
	"sort"
)

// Slice sorts the slice pointed by the provided pointer given the provided
// less function and removes repeated elements.
// The function panics if the provided interface is not a pointer to a slice.
func Slice(slicePtr interface{}, less func(i, j int) bool) {
	v := reflect.ValueOf(slicePtr).Elem()
	if v.Len() <= 1 {
		return
	}
	sort.Slice(v.Interface(), less)

	i := 0
	for j := 1; j < v.Len(); j++ {
		if !less(i, j) {
			continue
		}
		i++
		v.Index(i).Set(v.Index(j))
	}
	i++
	v.SetLen(i)
}
