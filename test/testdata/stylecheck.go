//args: -Estylecheck
package testdata

func Stylecheck(x int) {
	if 0 == x { // ERROR "don't use Yoda conditions"
		panic(x)
	}
}

func StylecheckNolintStylecheck(x int) {
	if 0 == x { //nolint:stylecheck
		panic(x)
	}
}

func StylecheckNolintMegacheck(x int) {
	if 0 == x { //nolint:megacheck // ERROR "don't use Yoda conditions"
		panic(x)
	}
}
