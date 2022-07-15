//golangcitest:args -Egofmt,goimports
package p

import (
    "os"
    "fmt"
)

 func goimports(a, b int) int {
 	if a != b {
 		return 1 
	}
 	return 2
}
