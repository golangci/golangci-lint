//args: -Estaticcheck
package testdata

func Staticcheck() {
	var x int
	x = x // ERROR "self-assignment of x to x"
}

func StaticcheckNolintStaticcheck() {
	var x int
	x = x //nolint:staticcheck
}

func StaticcheckNolintMegacheck() {
	var x int
	x = x //nolint:megacheck
}
