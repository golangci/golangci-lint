//args: -Estaticcheck
package testdata

import (
	"fmt"
	"runtime"
)

func Staticcheck() {
	var x int
	x = x // ERROR "self-assignment of x to x"
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

func StaticcheckDeprecated() {
	_ = runtime.CPUProfile() // ERROR "SA1019: runtime.CPUProfile has been deprecated .*"
}

func StaticcheckPrintf() {
	x := "dummy"
	fmt.Printf("%d", x) // ERROR "SA5009: Printf format %d has arg #1 of wrong type"
}
