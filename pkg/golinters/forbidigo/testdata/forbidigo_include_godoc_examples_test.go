//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/forbidigo_include_godoc_examples.yml
package testdata

import "fmt"

func ExampleForbidigoNoGodoc() {
	fmt.Printf("too noisy!!!") // want "use of `fmt.Printf` forbidden by pattern.*"
}
