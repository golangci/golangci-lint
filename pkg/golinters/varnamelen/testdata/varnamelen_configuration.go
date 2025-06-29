//golangcitest:args -Evarnamelen
//golangcitest:config_path testdata/varnamelen_configuration.yml
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func _() {
	ok := foo()

	fmt.Println("a")
	fmt.Println("b")
	fmt.Println("c")
	fmt.Println("d")
	println(ok)
}

func _() {
	fn := foo()

	fmt.Println("a")
	fmt.Println("b")
	fmt.Println("c")
	fmt.Println("d")
	println(fn)
}

func foo() bool {
	return true
}
