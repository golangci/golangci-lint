// args: -Etenv
package testdata

import (
	"os"
	"testing"
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

func TestF(t *testing.T) {
	os.Setenv("a", "b")                         // ERROR "func TestF is not using testing.Setenv"
	if err := os.Setenv("a", "b"); err != nil { // ERROR "func TestF is not using testing.Setenv"
		_ = err
	}
}

func BenchmarkF(b *testing.B) {
	os.Setenv("a", "b")                         // ERROR "func BenchmarkF is not using testing.Setenv"
	if err := os.Setenv("a", "b"); err != nil { // ERROR "func BenchmarkF is not using testing.Setenv"
		_ = err
	}
}

func testTB(tb testing.TB) {
	os.Setenv("a", "b")                         // ERROR "func testTB is not using testing.Setenv"
	if err := os.Setenv("a", "b"); err != nil { // ERROR "func testTB is not using testing.Setenv"
		_ = err
	}
}
