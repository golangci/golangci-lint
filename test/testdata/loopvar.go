//golangcitest:args -Eloopvar
package testdata

import (
	"fmt"
)

func foo() {
	slice := []int{1, 2, 3}
	fns := make([]func(), 0, len(slice)*2)
	for i, v := range slice {
		i := i // want `The loop variable "i" should not be copied \(GO_VERSION >= 1.22 or GOEXPERIMENT=loopvar\)`
		fns = append(fns, func() {
			fmt.Print(i)
		})
		_v := v // want `The loop variable "v" should not be copied \(GO_VERSION >= 1.22 or GOEXPERIMENT=loopvar\)`
		fns = append(fns, func() {
			fmt.Print(_v)
		})
	}
	for _, fn := range fns {
		fn()
	}
}

func bar() {
	loopCount := 3
	fns := make([]func(), 0, loopCount)
	for i := 1; i <= loopCount; i++ {
		i := i // want `The loop variable "i" should not be copied \(GO_VERSION >= 1.22 or GOEXPERIMENT=loopvar\)`
		fns = append(fns, func() {
			fmt.Print(i)
		})
	}
	for _, fn := range fns {
		fn()
	}
}
