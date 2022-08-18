//golangcitest:args -Egofmt,goimports
//golangcitest:expected_exitcode 0
package p

func goimports(a, b int) int {
	if a != b {
		return 1
	}
	return 2
}
