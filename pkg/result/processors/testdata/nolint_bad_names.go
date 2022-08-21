package testdata

import "math"

func RetErr() error {
	return nil
}

func MissedErrorCheck() {
	RetErr() //nolint:bad1,errcheck
}

//nolint:bad2,errcheck
func MissedErrorCheck2() {
	RetErr()
}

func _() {
	x := math.MinInt8
	for {
		_ = x
		x = 0 //nolint:bad1,ineffassign
		x = 0
	}
}
