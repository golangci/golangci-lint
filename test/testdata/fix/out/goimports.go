//args: -Egofmt,goimports
package p

func goimports(a, b int) int {
	if a != b {
		return 1
	}
	return 2
}
