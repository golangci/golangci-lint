package testdata

var y int

var unused int // nolint:megacheck // ERROR "`unused` is unused"

func f(x int) {
}

func g(x int) { // nolint:megacheck // ERROR "`g` is unused"
}

func H(x int) {
}

func init() {
	f(y)
}

var _ int
