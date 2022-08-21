//golangcitest:args -Edogsled
package testdata

func Dogsled() {
	_ = ret1()
	_, _ = ret2()
	_, _, _ = ret3()    // want "declaration has 3 blank identifiers"
	_, _, _, _ = ret4() // want "declaration has 4 blank identifiers"
}

func ret1() (a int) {
	return 1
}

func ret2() (a, b int) {
	return 1, 2
}

func ret3() (a, b, c int) {
	return 1, 2, 3
}

func ret4() (a, b, c, d int) {
	return 1, 2, 3, 4
}
