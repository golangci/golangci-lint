//golangcitest:args -Enlreturn
//golangcitest:config_path testdata/configs/nlreturn.yml
package testdata

func foo0(n int) int {
	if n == 1 {
		n2 := n * n
		return n2
	}

	return 1
}

func foo1(n int) int {
	if n == 1 {
		n2 := n * n
		n3 := n2 * n
		return n3 // want "return with no blank line before"
	}

	return 1
}
