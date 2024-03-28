//golangcitest:args -Euntypedconst
package testdata

import (
	"fmt"
)

type ExString string

func retExString() ExString {
	if true {
		return ExString("hoge")
	} else {
		fmt.Println("This should never happen")
		return "hoge" // want `returning untyped constant as defined type "command-line-arguments.ExString"`
	}
}

type ExInt int

func retExInt() ExInt {
	if true {
		return ExInt(1)
	} else {
		return 1 // want `returning untyped constant as defined type "command-line-arguments.ExInt"`
	}
}

type ExFloat float64

func retExFloat() ExFloat {
	if true {
		return ExFloat(0.5)
	} else {
		return 0.5 // want `returning untyped constant as defined type "command-line-arguments.ExFloat"`
	}
}

type ExComplex complex128

func retExComplex() ExComplex {
	if true {
		return ExComplex(1.0 + 0.5i)
	} else {
		return 1.0 + 0.5i // want `returning untyped constant as defined type "command-line-arguments.ExComplex"`
	}
}

type ExRune rune

func retExRune() ExRune {
	if true {
		return ExRune('a')
	} else {
		return 'a' // want `returning untyped constant as defined type "command-line-arguments.ExRune"`
	}
}

type ExBool bool

func retExBool() ExBool {
	if true {
		return ExBool(true)
	} else {
		return true // want `returning untyped constant as defined type "command-line-arguments.ExBool"`
	}
}
