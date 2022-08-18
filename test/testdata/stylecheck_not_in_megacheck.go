//golangcitest:args -Emegacheck
//golangcitest:expected_exitcode 0
package testdata

func StylecheckNotInMegacheck(x int) {
	if 0 == x {
		panic(x)
	}
}
