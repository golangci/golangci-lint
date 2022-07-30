//golangcitest:args -Enoloopclosure
package testdata

import _ "fmt"

func noloopclosureForLoop() {
	for i := 0; i < 5; i++ {
		_ = i
	}

	for false {
	}

	for i := 0; i < 5; i++ {
		_ = func() {
			_ = i // ERROR "found reference to loop variable `i`. Consider to duplicate variable `i` before using it inside the function closure."
		}
	}

	k := 5
	for i, j := 0, 0; i < j; i++ {
		_ = func() {
			_ = k
		}

		_ = func() {
			_, _ = i, j // ERROR "found reference to loop variable `i`. Consider to duplicate variable `i` before using it inside the function closure."
		}
	}
}

func noloopclosureRangeLoop() {
	for k, v := range map[string]int{} {
		_ = func() {
			_ = k // ERROR "found reference to loop variable `k`. Consider to duplicate variable `k` before using it inside the function closure."
			_ = v // ERROR "found reference to loop variable `v`. Consider to duplicate variable `v` before using it inside the function closure."
		}
	}

	for _, v := range map[string]int{} {
		_ = func() {
			_ = v // ERROR "found reference to loop variable `v`. Consider to duplicate variable `v` before using it inside the function closure."
		}
	}

	for k := range map[string]int{} {
		_ = func() {
			_ = k // ERROR "found reference to loop variable `k`. Consider to duplicate variable `k` before using it inside the function closure."
		}
	}

	for range map[string]int{} {
	}
}
