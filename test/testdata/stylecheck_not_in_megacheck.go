//args: -Emegacheck
package testdata

func StylecheckNotInMegacheck(x int) {
	if 0 == x {
		panic(x)
	}
}

func Staticcheck2() {
	var x int
	x = x // ERROR "self-assignment of x to x"
}
