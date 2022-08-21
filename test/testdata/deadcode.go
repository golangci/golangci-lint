//golangcitest:args -Edeadcode --internal-cmd-test
package testdata

var y int

var unused int // want "`unused` is unused"

func f(x int) {
}

func g(x int) { // want "`g` is unused"
}

func H(x int) {
}

func init() {
	f(y)
}

var _ int
