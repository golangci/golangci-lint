//golangcitest:args -Estylecheck
// Package testdata ...
package testdata

func Stylecheck(x int) {
	switch x {
	case 1:
		return
	default: // want "ST1015: default case should be first or last in switch statement"
		return
	case 2:
		return
	}
}

func StylecheckNolintStylecheck(x int) {
	switch x {
	case 1:
		return
	default: //nolint:stylecheck
		return
	case 2:
		return
	}
}

func StylecheckNolintMegacheck(x int) {
	switch x {
	case 1:
		return
	default: //nolint:megacheck // want "ST1015: default case should be first or last in switch statement"
		return
	case 2:
		return
	}
}
