//golangcitest:args -Ecopyloopvar
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func _() {
	slice := []int{1, 2, 3}
	fns := make([]func(), 0, len(slice)*2)
	for i, v := range slice {
		// want `The copy of the 'for' variable "i" can be deleted \(Go 1\.22\+\)`
		fns = append(fns, func() {
			fmt.Println(i)
		})
		// want `The copy of the 'for' variable "v" can be deleted \(Go 1\.22\+\)`
		fns = append(fns, func() {
			fmt.Println(v)
		})
		_v := v
		_ = _v
	}
	for _, fn := range fns {
		fn()
	}
}

func _() {
	loopCount := 3
	fns := make([]func(), 0, loopCount)
	for i := 1; i <= loopCount; i++ {
		// want `The copy of the 'for' variable "i" can be deleted \(Go 1\.22\+\)`
		fns = append(fns, func() {
			fmt.Println(i)
		})
	}
	for _, fn := range fns {
		fn()
	}
}
