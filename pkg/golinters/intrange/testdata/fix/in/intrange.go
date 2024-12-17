//go:build go1.22

//golangcitest:args -Eintrange
//golangcitest:expected_exitcode 0
package testdata

import "math"

func CheckIntrange() {
	for i := 0; i < 10; i++ {
	}

	for i := uint8(0); i < math.MaxInt8; i++ {
	}

	for i := 0; i < 10; i += 2 {
	}

	for i := 0; i < 10; i++ {
		i += 1
	}
}
