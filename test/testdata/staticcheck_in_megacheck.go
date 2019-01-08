//args: -Emegacheck
package testdata

func StaticcheckInMegacheck() {
	var x int
	x = x // ERROR "self-assignment of x to x"
}

func StaticcheckNolintStaticcheckInMegacheck() {
	var x int
	x = x //nolint:staticcheck
}

func StaticcheckNolintMegacheckInMegacheck() {
	var x int
	x = x //nolint:megacheck
}
