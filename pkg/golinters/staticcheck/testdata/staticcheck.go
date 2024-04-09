//golangcitest:args -Estaticcheck
package testdata

import (
	"fmt"
)

func Staticcheck() {
	var x int
	x = x // want "self-assignment of x to x"
	fmt.Printf("%d", x)
}

func StaticcheckNolintStaticcheck() {
	var x int
	x = x //nolint:staticcheck
}

func StaticcheckNolintMegacheck() {
	var x int
	x = x //nolint:megacheck
}

func StaticcheckPrintf() {
	x := "dummy"
	fmt.Printf("%d", x) // want "SA5009: Printf format %d has arg #1 of wrong type"
}
