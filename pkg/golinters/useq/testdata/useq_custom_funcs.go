//golangcitest:args -Euseq
//golangcitest:config_path testdata/useq_custom_funcs.yml
package testdata

import (
	"fmt"
)

func customPrint(msg string, args ...any) {
	fmt.Println(msg, args)
}

func main() {
	name := "John"

	customPrint("This is my name: \"%s\"\n", name) // want `use %q instead of \"%s\" for formatting strings with quotations`
	customPrint("This is my name: %s", name)
	customPrint("This is my name: %q", name)
	customPrint("This is my age: %d", 42)
}
