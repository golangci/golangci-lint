//golangcitest:args -Etenv
package testdata

import (
	"os"
	"testing"
)

var (
	e = os.Setenv("a", "b") // never seen
)

func setup() {
	os.Setenv("a", "b")        // OK
	err := os.Setenv("a", "b") // OK
	_ = err
	if err := os.Setenv("a", "b"); err != nil { // OK
		_ = err
	}
}

func TestF(t *testing.T) {
	os.Setenv("a", "b")        // want "os\\.Setenv\\(\\) can be replaced by `t\\.Setenv\\(\\)` in TestF"
	err := os.Setenv("a", "b") // want "os\\.Setenv\\(\\) can be replaced by `t\\.Setenv\\(\\)` in TestF"
	_ = err
	if err := os.Setenv("a", "b"); err != nil { // want "os\\.Setenv\\(\\) can be replaced by `t\\.Setenv\\(\\)` in TestF"
		_ = err
	}
}

func BenchmarkF(b *testing.B) {
	os.Setenv("a", "b")        // want "os\\.Setenv\\(\\) can be replaced by `b\\.Setenv\\(\\)` in BenchmarkF"
	err := os.Setenv("a", "b") // want "os\\.Setenv\\(\\) can be replaced by `b\\.Setenv\\(\\)` in BenchmarkF"
	_ = err
	if err := os.Setenv("a", "b"); err != nil { // want "os\\.Setenv\\(\\) can be replaced by `b\\.Setenv\\(\\)` in BenchmarkF"
		_ = err
	}
}

func testTB(tb testing.TB) {
	os.Setenv("a", "b")        // want "os\\.Setenv\\(\\) can be replaced by `tb\\.Setenv\\(\\)` in testTB"
	err := os.Setenv("a", "b") // want "os\\.Setenv\\(\\) can be replaced by `tb\\.Setenv\\(\\)` in testTB"
	_ = err
	if err := os.Setenv("a", "b"); err != nil { // want "os\\.Setenv\\(\\) can be replaced by `tb\\.Setenv\\(\\)` in testTB"
		_ = err
	}
}
