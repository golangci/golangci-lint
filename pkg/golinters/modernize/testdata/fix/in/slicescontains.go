//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package slicescontains

import "slices"

var _ = slices.Contains[[]int] // force import of "slices" to avoid duplicate import edits

func nopeNoBreak(slice []int, needle int) {
	for i := range slice {
		if slice[i] == needle {
			println("found")
		}
	}
}

func rangeIndex(slice []int, needle int) {
	for i := range slice { // want "Loop can be simplified using slices.Contains"
		if slice[i] == needle {
			println("found")
			break
		}
	}
}

func rangeValue(slice []int, needle int) {
	for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
		if elem == needle {
			println("found")
			break
		}
	}
}

func returns(slice []int, needle int) {
	for i := range slice { // want "Loop can be simplified using slices.Contains"
		if slice[i] == needle {
			println("found")
			return
		}
	}
}

func assignTrueBreak(slice []int, needle int) {
	found := false
	for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
		if elem == needle {
			found = true
			break
		}
	}
	print(found)
}

func assignFalseBreak(slice []int, needle int) {
	found := true
	for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
		if elem == needle {
			found = false
			break
		}
	}
	print(found)
}

func assignFalseBreakInSelectSwitch(slice []int, needle int) {
	// Exercise RangeStmt in CommClause, CaseClause.
	select {
	default:
		found := false
		for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
			if elem == needle {
				found = true
				break
			}
		}
		print(found)
	}
	switch {
	default:
		found := false
		for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
			if elem == needle {
				found = true
				break
			}
		}
		print(found)
	}
}

func returnTrue(slice []int, needle int) bool {
	for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
		if elem == needle {
			return true
		}
	}
	return false
}

func returnFalse(slice []int, needle int) bool {
	for _, elem := range slice { // want "Loop can be simplified using slices.Contains"
		if elem == needle {
			return false
		}
	}
	return true
}

func containsFunc(slice []int, needle int) bool {
	for _, elem := range slice { // want "Loop can be simplified using slices.ContainsFunc"
		if predicate(elem) {
			return true
		}
	}
	return false
}

func nopeLoopBodyHasFreeContinuation(slice []int, needle int) bool {
	for _, elem := range slice {
		if predicate(elem) {
			if needle == 7 {
				continue // this statement defeats loop elimination
			}
			return true
		}
	}
	return false
}

func generic[T any](slice []T, f func(T) bool) bool {
	for _, elem := range slice { // want "Loop can be simplified using slices.ContainsFunc"
		if f(elem) {
			return true
		}
	}
	return false
}

func predicate(int) bool

// Regression tests for bad fixes when needle
// and haystack have different types (#71313):

func nopeNeedleHaystackDifferentTypes(x any, args []error) {
	for _, arg := range args {
		if arg == x {
			return
		}
	}
}

func nopeNeedleHaystackDifferentTypes2(x error, args []any) {
	for _, arg := range args {
		if arg == x {
			return
		}
	}
}

func nopeVariadicNamedContainsFunc(slice []int) bool {
	for _, elem := range slice {
		if variadicPredicate(elem) {
			return true
		}
	}
	return false
}

func variadicPredicate(int, ...any) bool

func nopeVariadicContainsFunc(slice []int) bool {
	f := func(int, ...any) bool {
		return true
	}
	for _, elem := range slice {
		if f(elem) {
			return true
		}
	}
	return false
}

// Negative test case for implicit C->I conversion
type I interface{ F() }
type C int

func (C) F() {}

func nopeImplicitConversionContainsFunc(slice []C, f func(I) bool) bool {
	for _, elem := range slice {
		if f(elem) { // implicit conversion from C to I
			return true
		}
	}
	return false
}

func nopeTypeParamWidening[T any](slice []T, f func(any) bool) bool {
	for _, elem := range slice {
		if f(elem) { // implicit conversion from T to any
			return true
		}
	}
	return false
}
