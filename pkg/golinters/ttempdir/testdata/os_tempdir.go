//golangcitest:args -Ettempdir
package testdata

import (
	"os"
	"testing"
)

var (
	dir = os.TempDir() // never seen
)

func setup() {
	os.TempDir()        // never seen
	dir := os.TempDir() // never seen
	_ = dir
	_ = os.TempDir() // never seen
}

func F(t *testing.T) {
	setup()
	os.TempDir()                        // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	t.Log(os.TempDir())                 // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	_ = os.TempDir()                    // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	if dir := os.TempDir(); dir != "" { // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
		_ = dir
	}
}

func BF(b *testing.B) {
	TBF(b)
	os.TempDir()                        // want "os\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_ = os.TempDir()                    // want "os\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	if dir := os.TempDir(); dir != "" { // want "os\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
		_ = dir
	}
}

func TBF(tb testing.TB) {
	os.TempDir()                        // want "os\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_ = os.TempDir()                    // want "os\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	if dir := os.TempDir(); dir != "" { // want "os\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
		_ = dir
	}
}

func FF(f *testing.F) {
	os.TempDir()                        // want "os\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_ = os.TempDir()                    // want "os\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	if dir := os.TempDir(); dir != "" { // want "os\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
		_ = dir
	}
}
