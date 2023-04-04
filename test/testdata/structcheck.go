//golangcitest:args -Estructcheck
package testdata

type t struct {
	unusedField int // want "`unusedField` is unused"
}
