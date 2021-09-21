// args: -Etenv
// config_path: testdata/configs/tenv.yml
package testdata

import (
	"os"
)

var (
	e = os.Setenv("a", "b") // ERROR "variable e is not using t.Setenv"
)

func setup() {
	os.Setenv("a", "b")        // ERROR "func setup is not using t.Setenv"
	err := os.Setenv("a", "b") // ERROR "func setup is not using t.Setenv"
	if err != nil {
		_ = err
	}
}

func TestF() {
	os.Setenv("a", "b")                         // ERROR "func TestF is not using t.Setenv"
	if err := os.Setenv("a", "b"); err != nil { // ERROR "func TestF is not using t.Setenv"
		_ = err
	}
}
