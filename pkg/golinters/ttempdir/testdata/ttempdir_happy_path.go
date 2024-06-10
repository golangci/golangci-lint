//golangcitest:args -Ettempdir
//golangcitest:expected_exitcode 0
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
	t.TempDir()                        // never seen
	t.Log(t.TempDir())                 // never seen
	_ = t.TempDir()                    // never seen
	if dir := t.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func BF(b *testing.B) {
	TBF(b)
	b.TempDir()                        // never seen
	_ = b.TempDir()                    // never seen
	if dir := b.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func TBF(tb testing.TB) {
	tb.TempDir()                        // never seen
	_ = tb.TempDir()                    // never seen
	if dir := tb.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func FF(f *testing.F) {
	f.TempDir()                        // never seen
	_ = f.TempDir()                    // never seen
	if dir := f.TempDir(); dir != "" { // never seen
		_ = dir
	}
}
