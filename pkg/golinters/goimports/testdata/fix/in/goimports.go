//golangcitest:config_path testdata/goimports.yml
//golangcitest:expected_exitcode 0
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
