//go:build go1.22

//golangcitest:args -Eintrange
//golangcitest:expected_exitcode 0
package testdata

import "math"

func CheckIntrange() {
	for range 10 {
	}

	for range uint8(math.MaxInt8) {
	}

	for i := 0; i < 10; i += 2 {
	}

	for i := 0; i < 10; i++ {
		i += 1
	}
}
