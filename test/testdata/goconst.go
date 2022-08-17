//golangcitest:args -Egoconst
package testdata

import "fmt"

func GoconstA() {
	a := "needconst" // want "string `needconst` has 5 occurrences, make it a constant"
	fmt.Print(a)
	b := "needconst"
	fmt.Print(b)
	c := "needconst"
	fmt.Print(c)
}

func GoconstB() {
	a := "needconst"
	fmt.Print(a)
	b := "needconst"
	fmt.Print(b)
}

const AlreadyHasConst = "alreadyhasconst"

func GoconstC() {
	a := "alreadyhasconst" // want "string `alreadyhasconst` has 3 occurrences, but such constant `AlreadyHasConst` already exists"
	fmt.Print(a)
	b := "alreadyhasconst"
	fmt.Print(b)
	c := "alreadyhasconst"
	fmt.Print(c)
	fmt.Print("alreadyhasconst")
}
