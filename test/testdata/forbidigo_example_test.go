//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/configs/forbidigo.yml
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func ExampleForbidigo() {
	fmt.Printf("too noisy!!!") // godoc examples are ignored (in *_test.go files only)
}
