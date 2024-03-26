//golangcitest:args -Etenv
package testdata

import (
	"os"
	"testing"
)

func FuzzF(f *testing.F) {
	os.Setenv("a", "b")        // want "os\\.Setenv\\(\\) can be replaced by `f\\.Setenv\\(\\)` in FuzzF"
	err := os.Setenv("a", "b") // want "os\\.Setenv\\(\\) can be replaced by `f\\.Setenv\\(\\)` in FuzzF"
	_ = err
	if err := os.Setenv("a", "b"); err != nil { // want "os\\.Setenv\\(\\) can be replaced by `f\\.Setenv\\(\\)` in FuzzF"
		_ = err
	}
}
