//golangcitest:args -Ettempdir
//golangcitest:config_path testdata/ttempdir_all.yml
package testdata

import (
	"os"
	"testing"
)

var (
	_, e = os.MkdirTemp("a", "b") // never seen
)

func testsetup() {
	os.MkdirTemp("a", "b")           // if -all = true, want  "os\\.MkdirTemp\\(\\) should be replaced by `testing\\.TempDir\\(\\)` in testsetup"
	_, err := os.MkdirTemp("a", "b") // if -all = true, want  "os\\.MkdirTemp\\(\\) should be replaced by `testing\\.TempDir\\(\\)` in testsetup"
	if err != nil {
		_ = err
	}
	os.MkdirTemp("a", "b") // if -all = true, "func setup is not using testing.TempDir"
}

func TestF(t *testing.T) {
	testsetup()
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestF"
		_ = err
	}
}

func BenchmarkF(b *testing.B) {
	TB(b)
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BenchmarkF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BenchmarkF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BenchmarkF"
		_ = err
	}
}

func TB(tb testing.TB) {
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TB"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TB"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TB"
		_ = err
	}
}

func FuzzF(f *testing.F) {
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FuzzF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FuzzF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FuzzF"
		_ = err
	}
}

func TestFunctionLiteral(t *testing.T) {
	testsetup()
	t.Run("test", func(t *testing.T) {
		os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
		_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
		_ = err
		if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
			_ = err
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Run("test", func(*testing.T) {})
}

func TestEmptyTB(t *testing.T) {
	func(testing.TB) {}(t)
}
