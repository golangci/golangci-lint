//go:build go1.22

//golangcitest:args -Ecopyloopvar
package testdata

import "fmt"

func copyloopvarCase1() {
	slice := []int{1, 2, 3}
	fns := make([]func(), 0, len(slice)*2)
	for i, v := range slice {
		i := i // want `The copy of the 'for' variable "i" can be deleted \(Go 1\.22\+\)`
		fns = append(fns, func() {
			fmt.Println(i)
		})
		_v := v // want `The copy of the 'for' variable "v" can be deleted \(Go 1\.22\+\)`
		fns = append(fns, func() {
			fmt.Println(_v)
		})
	}
	for _, fn := range fns {
		fn()
	}
}

func copyloopvarCase2() {
	loopCount := 3
	fns := make([]func(), 0, loopCount)
	for i := 1; i <= loopCount; i++ {
		i := i // want `The copy of the 'for' variable "i" can be deleted \(Go 1\.22\+\)`
		fns = append(fns, func() {
			fmt.Println(i)
		})
	}
	for _, fn := range fns {
		fn()
	}
}
