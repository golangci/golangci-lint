//golangcitest:args -Efoo --simple --hello=world
//golangcitest:config_path testdata/example.yml
//golangcitest:expected_linter bar
package testdata

import "fmt"

func main() {
	fmt.Println("Hello")
}
