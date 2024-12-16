//golangcitest:args -Estaticcheck
//golangcitest:expected_exitcode 0
package testdata

import "sort"

func _(a []string) {
	a = sort.StringSlice(a)
}
