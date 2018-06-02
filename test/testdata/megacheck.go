package testdata

func Megacheck() {
	var x int
	x = x // nolint:ineffassign // ERROR "self-assignment of x to x"
}
