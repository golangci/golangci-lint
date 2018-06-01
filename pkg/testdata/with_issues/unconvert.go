package testdata

import "log" // nolint:depguard

func Unconvert() {
	a := 1
	b := int(a) // ERROR "unnecessary conversion"
	log.Print(b)
}
