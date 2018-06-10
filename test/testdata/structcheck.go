// args: -Estructcheck
package testdata

type t struct {
	unusedField int // ERROR "`unusedField` is unused"
}
