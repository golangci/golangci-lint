//args: -Edeadcode
package testdata

var y int

var unused int // ERROR "`unused` is unused"

func f(x int) {
}

func g(x int) { // ERROR "`g` is unused"
}

func H(x int) {
}

func init() {
	f(y)
}

var _ int
