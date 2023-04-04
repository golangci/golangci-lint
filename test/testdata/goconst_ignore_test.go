//golangcitest:args -Egoconst
//golangcitest:config_path testdata/configs/goconst_ignore.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"testing"
)

func TestGoConstA(t *testing.T) {
	a := "needconst"
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
	a := "alreadyhasconst"
	fmt.Print(a)
	b := "alreadyhasconst"
	fmt.Print(b)
	c := "alreadyhasconst"
	fmt.Print(c)
	fmt.Print("alreadyhasconst")
}
