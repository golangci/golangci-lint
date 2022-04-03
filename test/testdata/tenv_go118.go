//go:build go1.18
// +build go1.18

// args: -Etenv
package testdata

import (
	"os"
	"testing"
)

func FuzzF(f *testing.F) {
	os.Setenv("a", "b")        // ERROR "os\\.Setenv\\(\\) can be replaced by `f\\.Setenv\\(\\)` in FuzzF"
	err := os.Setenv("a", "b") // ERROR "os\\.Setenv\\(\\) can be replaced by `f\\.Setenv\\(\\)` in FuzzF"
	_ = err
	if err := os.Setenv("a", "b"); err != nil { // ERROR "os\\.Setenv\\(\\) can be replaced by `f\\.Setenv\\(\\)` in FuzzF"
		_ = err
	}
}
