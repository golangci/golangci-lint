//golangcitest:args -Enoloopclosure
package testdata

import _ "fmt"

func noloopclosureForLoop() {
	for {

	}

	for false {
		// noop
	}

	for i := 0; ; i++ {
		_ = func() {
			_ = i // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	for i := 0; ; {
		_ = i
	}

	for i := 0; i < 5; i++ {
		_ = i
	}

	var i int
	for i < 5 {
		// Note: we ignore the condition clause to reduce false positive.
		// Because it's hard to check which variable inside the condition that will be mutated.
		_ = func() {
			_ = i
		}
	}

	for i := 0; i < 5; i++ {
		_ = func() {
			_ = i // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			_ = func() {
				_ = i // ERROR "found captured reference to loop variable inside a closure"
				_ = j // ERROR "found captured reference to loop variable inside a closure"
			}
		}

		var j int
		_ = func() {
			_ = j
		}
	}

	k := 5
	for i, j := 0, 0; i < j; i, k = i+1, k+1 {
		_ = func() {
			// Not okay since it's listed in the PostInit.
			_ = k // ERROR "found captured reference to loop variable inside a closure"
		}

		_ = func() {
			_ = i // ERROR "found captured reference to loop variable inside a closure"

			// Not okay, even if it's not part of the PostInit, as it's very likely that dev will mutate it somehow
			// inside the for loop.
			_ = j // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	// Can handle ast.SelectorExpr properly.
	x := struct {
		A int
		B int
	}{}
	for x.A = 0; x.A < 5; x.A++ {
		_ = func() {
			_ = x.A // ERROR "found captured reference to loop variable inside a closure"
			_ = x.B
		}
	}

	// Can handle ast.ParenExpr properly.
	for (k) = 0; k < 5; k++ {
		_ = func() {
			_ = k // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	// Can handle ast.StarExpr properly.
	var p *int = new(int)
	for *p = 0; *p < 5; (*p)++ {
		_ = func() {
			_ = *p // ERROR "found captured reference to loop variable inside a closure"
			_ = p  // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	// Can handle ast.IndexExpr properly.
	arrayOfInt := []int{1, 2, 3}
	for arrayOfInt[1] = 0; arrayOfInt[1] < 5; arrayOfInt[1]++ {
		_ = func() {
			// Note: we will disallow any captured reference to the array, not just arrayOfInt[1]
			// because we cant accurately know which index that is being mutated in compile time.
			//
			// Generally it's a good practice to not do this anyway hence we raise it as an issue.
			_ = arrayOfInt[1] // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	mapOfInt := map[int]int{1: 1}
	for mapOfInt[1] = 0; mapOfInt[1] < 5; mapOfInt[1]++ {
		_ = func() {
			// Note: disallow any captured reference to the map.
			_ = mapOfInt[1] // ERROR "found captured reference to loop variable inside a closure"
		}
	}
}

func noloopclosureRangeLoop() {
	for k, v := range map[string]int{} {
		_ = func() {
			_ = k // ERROR "found captured reference to loop variable inside a closure"
			_ = v // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	for _, v := range map[string]int{} {
		_ = func() {
			_ = v // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	for k := range map[string]int{} {
		_ = func() {
			_ = k // ERROR "found captured reference to loop variable inside a closure"
		}
	}

	for range map[string]int{} {
	}

	x := struct{ A string }{}
	for x.A = range map[string]int{} {
		_ = func() {
			_ = x.A // ERROR "found captured reference to loop variable inside a closure"
		}
	}
}
