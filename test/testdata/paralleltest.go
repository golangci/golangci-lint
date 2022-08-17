//golangcitest:args -Eparalleltest
package testdata

import (
	"fmt"
	"testing"
)

func TestFunctionSuccessfulRangeTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionMissingCallToParallel(t *testing.T) {} // ERROR "Function TestFunctionMissingCallToParallel missing the call to method parallel"
