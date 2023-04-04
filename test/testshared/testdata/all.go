//golangcitest:args -Efoo --simple --hello=world
//golangcitest:config_path testdata/example.yml
//golangcitest:expected_linter bar
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

func main() {
	fmt.Println("Hello")
}
