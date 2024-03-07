//golangcitest:args -Egocritic,gofumpt
//golangcitest:config_path testdata/configs/multiple-issues-fix.yml
//golangcitest:expected_exitcode 0
package p

import "fmt"

func main() {
	//standard greeting
	fmt.Println("hello world")
}
