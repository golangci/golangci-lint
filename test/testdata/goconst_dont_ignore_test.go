//golangcitest:args -Egoconst
//golangcitest:config_path testdata/configs/goconst_dont_ignore.yml
package testdata

import (
	"fmt"
	"testing"
)

func TestGoConstA(t *testing.T) {
	a := "needconst" // want "string `needconst` has 5 occurrences, make it a constant"
	fmt.Print(a)
	b := "needconst"
	fmt.Print(b)
	c := "needconst"
	fmt.Print(c)
}

func TestGoConstB(t *testing.T) {
	a := "needconst"
	fmt.Print(a)
	b := "needconst"
	fmt.Print(b)
}

const AlreadyHasConst = "alreadyhasconst"

func TestGoConstC(t *testing.T) {
	a := "alreadyhasconst" // want "string `alreadyhasconst` has 3 occurrences, but such constant `AlreadyHasConst` already exists"
	fmt.Print(a)
	b := "alreadyhasconst"
	fmt.Print(b)
	c := "alreadyhasconst"
	fmt.Print(c)
	fmt.Print("alreadyhasconst")
}
