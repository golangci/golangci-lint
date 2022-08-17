//go:build go1.18
// +build go1.18

//golangcitest:args -Ethelper
package testdata

import "testing"

func fhelperWithHelperAfterAssignment(f *testing.F) { // ERROR "test helper function should start from f.Helper()"
	_ = 0
	f.Helper()
}

func fhelperWithNotFirst(s string, f *testing.F, i int) { // ERROR `parameter \*testing.F should be the first`
	f.Helper()
}

func fhelperWithIncorrectName(o *testing.F) { // ERROR `parameter \*testing.F should have name f`
	o.Helper()
}

func FuzzSubtestShouldNotBeChecked(f *testing.F) {
	f.Add(5, "hello")
	f.Fuzz(func(t *testing.T, a int, b string) {})
}
