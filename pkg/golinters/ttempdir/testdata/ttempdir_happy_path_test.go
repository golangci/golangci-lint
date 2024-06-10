//golangcitest:args -Ettempdir
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	tdir = os.TempDir() // never seen
)

func testsetup() {
	os.TempDir()        // never seen
	dir := os.TempDir() // never seen
	_ = dir
	_ = os.TempDir() // never seen
}

func TestF(t *testing.T) {
	testsetup()
	t.TempDir()                       // never seen
	_ = t.TempDir()                   // never seen
	if dir = t.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func BenchmarkF(b *testing.B) {
	TB(b)
	b.TempDir()                       // never seen
	_ = b.TempDir()                   // never seen
	if dir = b.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func TB(tb testing.TB) {
	tb.TempDir()                       // never seen
	_ = tb.TempDir()                   // never seen
	if dir = tb.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func FuzzF(f *testing.F) {
	f.TempDir()                       // never seen
	_ = f.TempDir()                   // never seen
	if dir = f.TempDir(); dir != "" { // never seen
		_ = dir
	}
}

func TestFunctionLiteral(t *testing.T) {
	testsetup()
	t.Run("test", func(t *testing.T) {
		t.TempDir()                       // never seen
		_ = t.TempDir()                   // never seen
		if dir = t.TempDir(); dir != "" { // never seen
			_ = dir
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Run("test", func(*testing.T) {})
}

func TestEmptyTB(t *testing.T) {
	func(testing.TB) {}(t)
}

func TestRecursive(t *testing.T) {
	t.Log( // recursion level 1
		t.TempDir(), // never seen
	)
	t.Log( // recursion level 1
		fmt.Sprintf("%s", // recursion level 2
			t.TempDir(), // never seen
		),
	)
	t.Log( // recursion level 1
		filepath.Clean( // recursion level 2
			fmt.Sprintf("%s", // recursion level 3
				t.TempDir(), // never seen
			),
		),
	)
	t.Log( // recursion level 1
		filepath.Join( // recursion level 2
			filepath.Clean( // recursion level 3
				fmt.Sprintf("%s", // recursion level 4
					t.TempDir(), // never seen
				),
			),
			"test",
		),
	)
	t.Log( // recursion level 1
		fmt.Sprintf("%s/foo-%d", // recursion level 2
			filepath.Join( // recursion level 3
				filepath.Clean( // recursion level 4
					fmt.Sprintf("%s", // recursion level 5
						os.TempDir(), // max recursion level reached.
					),
				),
				"test",
			),
			1024,
		),
	)
}
