// args: -Enosnakecase
package testdata

func a_() { // ERROR "a_ is used under score. You should use mixedCap or MixedCap."
}

func b(a_a int) { // ERROR "a_a is used under score. You should use mixedCap or MixedCap."
}

func c() (c_c int) { // ERROR "c_c is used under score. You should use mixedCap or MixedCap."
	c_c = 1    // ERROR "c_c is used under score. You should use mixedCap or MixedCap."
	return c_c // It's never detected, because `c_c` is already detected.
}

func d() {
	var d_d int // ERROR "d_d is used under score. You should use mixedCap or MixedCap."
	_ = d_d     // It's never detected, because `_` is meaningful in Go and `d_d` is already detected.
}
