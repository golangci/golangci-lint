//golangcitest:args -Emegacheck
package testdata

import "fmt"

func StaticcheckInMegacheck() {
	var x int
	x = x // ERROR staticcheck "self-assignment of x to x"
	fmt.Printf("%d", x)
}

func StaticcheckNolintStaticcheckInMegacheck() {
	var x int
	x = x //nolint:staticcheck
}

func StaticcheckNolintMegacheckInMegacheck() {
	var x int
	x = x //nolint:megacheck
}

func Staticcheck2() {
	var x int
	x = x // ERROR staticcheck "self-assignment of x to x"
	fmt.Printf("%d", x)
}
