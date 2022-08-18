//golangcitest:args -Emegacheck
package testdata

func StylecheckNotInMegacheck(x int) {
	if 0 == x {
		panic(x)
	}
}
