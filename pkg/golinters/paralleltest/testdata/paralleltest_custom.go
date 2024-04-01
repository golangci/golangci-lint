//golangcitest:args -Eparalleltest
//golangcitest:config_path testdata/paralleltest_custom.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"testing"
)

func TestParallelTestIgnore(t *testing.T) {
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

func TestParallelTestEmptyIgnore(t *testing.T) {}
