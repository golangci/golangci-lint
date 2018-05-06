package testdata

type t struct { // ERROR "`t` is unused"
	unusedField int // ERROR "`testdata.t.unusedField` is unused"
}
