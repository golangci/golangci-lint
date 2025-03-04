//golangcitest:args -Estaticcheck
//golangcitest:config_path testdata/stylecheck.yml
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
	default: //nolint:staticcheck
		return
	case 2:
		return
	}
}
