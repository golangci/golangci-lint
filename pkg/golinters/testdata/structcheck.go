package testdata

type t struct { // ERROR "`t` is unused"
	unusedField int // ERROR "`unusedField` is unused"
}
