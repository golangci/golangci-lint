//golangcitest:args -Eexptostd
package testdata

import (
	"fmt"

	"golang.org/x/exp/maps"   // want `Import statement 'golang.org/x/exp/maps' can be replaced by 'maps'`
	"golang.org/x/exp/slices" // want `Import statement 'golang.org/x/exp/slices' can be replaced by 'slices'`
)

func _(m, a map[string]string) {
	maps.Clone(m) // want `golang.org/x/exp/maps.Clone\(\) can be replaced by maps.Clone\(\)`

	maps.Equal(m, a) // want `golang.org/x/exp/maps.Equal\(\) can be replaced by maps.Equal\(\)`

	maps.EqualFunc(m, a, func(i, j string) bool { // want `golang.org/x/exp/maps.EqualFunc\(\) can be replaced by maps.EqualFunc\(\)`
		return true
	})

	maps.Copy(m, a) // want `golang.org/x/exp/maps.Copy\(\) can be replaced by maps.Copy\(\)`

	maps.DeleteFunc(m, func(_, _ string) bool { // want `golang.org/x/exp/maps.DeleteFunc\(\) can be replaced by maps.DeleteFunc\(\)`
		return true
	})

	maps.Clear(m) // want `golang.org/x/exp/maps.Clear\(\) can be replaced by clear\(\)`

	fmt.Println("Hello")
}

func _(a, b []string) {
	slices.Equal(a, b)
	slices.EqualFunc(a, b, func(_ string, _ string) bool {
		return true
	})
	slices.Compare(a, b)
	slices.CompareFunc(a, b, func(_ string, _ string) int {
		return 0
	})
	slices.Index(a, "a")
	slices.IndexFunc(a, func(_ string) bool {
		return true
	})
	slices.Contains(a, "a")
	slices.ContainsFunc(a, func(_ string) bool {
		return true
	})
	slices.Insert(a, 0, "a", "b")
	slices.Delete(a, 0, 1)
	slices.DeleteFunc(a, func(_ string) bool {
		return true
	})
	slices.Replace(a, 0, 1, "a")
	slices.Clone(a)
	slices.Compact(a)
	slices.CompactFunc(a, func(_ string, _ string) bool {
		return true
	})
	slices.Grow(a, 2)
	slices.Clip(a)
	slices.Reverse(a)
	slices.Sort(a)
	slices.SortFunc(a, func(_, _ string) int {
		return 0
	})
	slices.SortStableFunc(a, func(_, _ string) int {
		return 0
	})
	slices.IsSorted(a)
	slices.IsSortedFunc(a, func(_, _ string) int {
		return 0
	})
	slices.Min(a)
	slices.MinFunc(a, func(_, _ string) int {
		return 0
	})
	slices.Max(a)
	slices.MaxFunc(a, func(_, _ string) int {
		return 0
	})
	slices.BinarySearch(a, "a")
	slices.BinarySearchFunc(a, b, func(_ string, _ []string) int {
		return 0
	})
}
