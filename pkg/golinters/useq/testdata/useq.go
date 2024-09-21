//golangcitest:args -Euseq
package testdata

import (
	"fmt"
	"os"
)

func main() {
	name := "John"

	fmt.Printf("This is my name: \"%s\"\n", name) // want `use %q instead of \"%s\" for formatting strings with quotations`
	fmt.Printf("This is my name: %s", name)
	fmt.Printf("This is my name: %q", name)
	fmt.Printf("This is my age: %d", 42)

	fmt.Sprintf("This is my name: \"%s\"\n", name) // want `use %q instead of \"%s\" for formatting strings with quotations`
	fmt.Sprintf("This is my name: %s", name)
	fmt.Sprintf("This is my name: %q", name)
	fmt.Sprintf("This is my age: %d", 42)

	fmt.Fprintf(os.Stdout, "This is my name: \"%s\"\n", name) // want `use %q instead of \"%s\" for formatting strings with quotations`
	fmt.Fprintf(os.Stdout, "This is my name: %s", name)
	fmt.Fprintf(os.Stdout, "This is my name: %q", name)
	fmt.Fprintf(os.Stdout, "This is my age: %d", 42)

	fmt.Errorf("This is my name: \"%s\"\n", name) // want `use %q instead of \"%s\" for formatting strings with quotations`
	fmt.Errorf("This is my name: %s", name)
	fmt.Errorf("This is my name: %q", name)
	fmt.Errorf("This is my age: %d", 42)
}
