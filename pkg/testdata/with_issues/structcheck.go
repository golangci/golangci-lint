package testdata

type t struct { // nolint:megacheck // ERROR "`t` is unused"
	unusedField int // nolint:megacheck // ERROR "`unusedField` is unused"
}
