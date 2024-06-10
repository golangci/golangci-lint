//golangcitest:args -Ettempdir
//golangcitest:config_path testdata/ttempdir_all.yml
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
}

func BF(b *testing.B) {
	TBF(b)
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
		_ = err
	}
}

func TBF(tb testing.TB) {
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
		_ = err
	}
}

func FF(f *testing.F) {
	os.MkdirTemp("a", "b")           // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_, err := os.MkdirTemp("a", "b") // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_ = err
	if _, err := os.MkdirTemp("a", "b"); err != nil { // want "os\\.MkdirTemp\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
		_ = err
	}
}
