//golangcitest:args -Egosimple
//golangcitest:expected_exitcode 0
package testdata

func _(src []string) {
	var dst []string
	copy(dst, src)
}
