//golangcitest:args -Eexptostd
package testdata

import (
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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
	slices.Equal(a, b) // want `golang.org/x/exp/slices\.Equal\(\) can be replaced by slices\.Equal\(\)`

	slices.EqualFunc(a, b, func(_ string, _ string) bool { // want `golang.org/x/exp/slices\.EqualFunc\(\) can be replaced by slices\.EqualFunc\(\)`
		return true
	})

	slices.Compare(a, b) // want `golang.org/x/exp/slices\.Compare\(\) can be replaced by slices\.Compare\(\)`

	slices.CompareFunc(a, b, func(_ string, _ string) int { // want `golang.org/x/exp/slices\.CompareFunc\(\) can be replaced by slices\.CompareFunc\(\)`
		return 0
	})

	slices.Index(a, "a") // want `golang.org/x/exp/slices\.Index\(\) can be replaced by slices\.Index\(\)`

	slices.IndexFunc(a, func(_ string) bool { // want `golang.org/x/exp/slices\.IndexFunc\(\) can be replaced by slices\.IndexFunc\(\)`
		return true
	})

	slices.Contains(a, "a") // want `golang.org/x/exp/slices\.Contains\(\) can be replaced by slices\.Contains\(\)`

	slices.ContainsFunc(a, func(_ string) bool { // want `golang.org/x/exp/slices\.ContainsFunc\(\) can be replaced by slices\.ContainsFunc\(\)`
		return true
	})

	slices.Insert(a, 0, "a", "b") // want `golang.org/x/exp/slices\.Insert\(\) can be replaced by slices\.Insert\(\)`

	slices.Delete(a, 0, 1) // want `golang.org/x/exp/slices\.Delete\(\) can be replaced by slices\.Delete\(\)`

	slices.DeleteFunc(a, func(_ string) bool { // want `golang.org/x/exp/slices\.DeleteFunc\(\) can be replaced by slices\.DeleteFunc\(\)`
		return true
	})

	slices.Replace(a, 0, 1, "a") // want `golang.org/x/exp/slices\.Replace\(\) can be replaced by slices\.Replace\(\)`

	slices.Clone(a) // want `golang.org/x/exp/slices\.Clone\(\) can be replaced by slices\.Clone\(\)`

	slices.Compact(a) // want `golang.org/x/exp/slices\.Compact\(\) can be replaced by slices\.Compact\(\)`

	slices.CompactFunc(a, func(_ string, _ string) bool { // want `golang.org/x/exp/slices\.CompactFunc\(\) can be replaced by slices\.CompactFunc\(\)`
		return true
	})

	slices.Grow(a, 2) // want `golang.org/x/exp/slices\.Grow\(\) can be replaced by slices\.Grow\(\)`

	slices.Clip(a) // want `golang.org/x/exp/slices\.Clip\(\) can be replaced by slices\.Clip\(\)`

	slices.Reverse(a) // want `golang.org/x/exp/slices\.Reverse\(\) can be replaced by slices\.Reverse\(\)`

	slices.Sort(a) // want `golang.org/x/exp/slices\.Sort\(\) can be replaced by slices\.Sort\(\)`

	slices.SortFunc(a, func(_, _ string) int { // want `golang.org/x/exp/slices\.SortFunc\(\) can be replaced by slices\.SortFunc\(\)`
		return 0
	})

	slices.SortStableFunc(a, func(_, _ string) int { // want `golang.org/x/exp/slices\.SortStableFunc\(\) can be replaced by slices\.SortStableFunc\(\)`
		return 0
	})

	slices.IsSorted(a) // want `golang.org/x/exp/slices\.IsSorted\(\) can be replaced by slices\.IsSorted\(\)`

	slices.IsSortedFunc(a, func(_, _ string) int { // want `golang.org/x/exp/slices\.IsSortedFunc\(\) can be replaced by slices\.IsSortedFunc\(\)`
		return 0
	})

	slices.Min(a) // want `golang.org/x/exp/slices\.Min\(\) can be replaced by slices\.Min\(\)`

	slices.MinFunc(a, func(_, _ string) int { // want `golang.org/x/exp/slices\.MinFunc\(\) can be replaced by slices\.MinFunc\(\)`
		return 0
	})

	slices.Max(a) // want `golang.org/x/exp/slices\.Max\(\) can be replaced by slices\.Max\(\)`

	slices.MaxFunc(a, func(_, _ string) int { // want `golang.org/x/exp/slices\.MaxFunc\(\) can be replaced by slices\.MaxFunc\(\)`
		return 0
	})

	slices.BinarySearch(a, "a") // want `golang.org/x/exp/slices\.BinarySearch\(\) can be replaced by slices\.BinarySearch\(\)`

	slices.BinarySearchFunc(a, b, func(_ string, _ []string) int { // want `golang.org/x/exp/slices\.BinarySearchFunc\(\) can be replaced by slices\.BinarySearchFunc\(\)`
		return 0
	})
}
