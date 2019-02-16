//args: -Egofmt
package p

func gofmt(a, b int) int {
	if a != b {
		return 1
	}
	return 2
}
