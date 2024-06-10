//golangcitest:args -Ettempdir
package testdata

import (
	"os"
	"testing"
)

var (
	_, ee = os.MkdirTemp("a", "b") // never seen
)

func setup() {
	os.MkdirTemp("a", "b")           // never seen
	_, err := os.MkdirTemp("a", "b") // never seen
	if err != nil {
		_ = err
	}
	os.MkdirTemp("a", "b") // never seen
}

func F(t *testing.T) {
	setup()
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
		_ = err
	}

	_ = func(t *testing.T) {
		_ = t
		_, _ = os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
	}

	t.Cleanup(func() {
		_, _ = os.MkdirTemp("a", "b")
	})
}

func BF(b *testing.B) {
	TBF(b)
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
		_ = err
	}

	func(b *testing.B) {
		_ = b
		_, _ = os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in anonymous function"
	}(b)
}

func TBF(tb testing.TB) {
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
		_ = err
	}

	defer func(tb testing.TB) {
		_ = tb
		_, _ = os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in anonymous function"
	}(tb)
}

func FF(f *testing.F) {
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
		_ = err
	}

	defer os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
}
