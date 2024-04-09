//golangcitest:args -Egofmt
//golangcitest:expected_exitcode 0
package p

 func gofmt(a, b int) int {
 	if a != b {
 		return 1 
	}
 	return 2
}