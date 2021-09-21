// args: -Etenv
package testdata

import (
	"os"
)

var (
	e = os.Setenv("a", "b") // OK
)

func setup() {
	os.Setenv("a", "b")        // OK
	err := os.Setenv("a", "b") // OK
	if err != nil {
		_ = err
	}
}

func TestF() {
	os.Setenv("a", "b")                         // OK
	if err := os.Setenv("a", "b"); err != nil { // OK
		_ = err
	}
}
