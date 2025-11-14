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
	if slices.Contains(slice, needle) {
		println("found")
	}
}

func rangeValue(slice []int, needle int) {
	if slices.Contains(slice, needle) {
		println("found")
	}
}

func returns(slice []int, needle int) {
	if slices.Contains(slice, needle) {
		println("found")
		return
	}
}

func assignTrueBreak(slice []int, needle int) {
	found := slices.Contains(slice, needle)
	print(found)
}

func assignFalseBreak(slice []int, needle int) {
	found := !slices.Contains(slice, needle)
	print(found)
}

func assignFalseBreakInSelectSwitch(slice []int, needle int) {
	// Exercise RangeStmt in CommClause, CaseClause.
	select {
	default:
		found := slices.Contains(slice, needle)
		print(found)
	}
	switch {
	default:
		found := slices.Contains(slice, needle)
		print(found)
	}
}

func returnTrue(slice []int, needle int) bool {
	return slices.Contains(slice, needle)
}

func returnFalse(slice []int, needle int) bool {
	return !slices.Contains(slice, needle)
}

func containsFunc(slice []int, needle int) bool {
	return slices.ContainsFunc(slice, predicate)
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
	return slices.ContainsFunc(slice, f)
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
