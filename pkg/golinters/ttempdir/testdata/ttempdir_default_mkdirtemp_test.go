//golangcitest:args -Ettempdir
package testdata

import (
	"fmt"
	"os"
	"testing"
)

var (
	_, e = os.MkdirTemp("a", "b") // never seen
)

func testsetup() {
	os.MkdirTemp("a", "b")           // never seen
	_, err := os.MkdirTemp("a", "b") // never seen
	if err != nil {
		_ = err
	}
	os.MkdirTemp("a", "b") // never seen
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

func TestTDD(t *testing.T) {
	for _, tt := range []struct {
		name string
	}{
		{"test"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
			_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
			_ = err
			if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
				_ = err
			}
		})
	}
}

func TestLoop(t *testing.T) {
	for i := 0; i < 3; i++ {
		os.MkdirTemp(fmt.Sprintf("a%d", i), "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestLoop"
		_, err := os.MkdirTemp(fmt.Sprintf("a%d", i), "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestLoop"
		_ = err
		if _, err := os.MkdirTemp(fmt.Sprintf("a%d", i), "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestLoop"
			_ = err
		}
	}
}
