//args: -Emegacheck
package testdata

func Megacheck() {
	var x int
	x = x // ERROR "self-assignment of x to x"
}
