//golangcitest:args -Ettempdir
package testdata

import (
	"io/ioutil"
	"testing"
)

var (
	_, e = ioutil.TempDir("a", "b") // never seen
)

func testsetup() {
	ioutil.TempDir("a", "b")           // if -all = true, want  "ioutil\\.TempDir\\(\\) should be replaced by `testing\\.TempDir\\(\\)` in testsetup"
	_, err := ioutil.TempDir("a", "b") // if -all = true, want  "ioutil\\.TempDir\\(\\) should be replaced by `testing\\.TempDir\\(\\)` in testsetup"
	if err != nil {
		_ = err
	}
	ioutil.TempDir("a", "b") // if -all = true, "func setup is not using testing.TempDir"
}

func TestF(t *testing.T) {
	testsetup()
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestF"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestF"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestF"
		_ = err
	}
}

func BenchmarkF(b *testing.B) {
	TB(b)
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BenchmarkF"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BenchmarkF"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BenchmarkF"
		_ = err
	}
}

func TB(tb testing.TB) {
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TB"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TB"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TB"
		_ = err
	}
}

func FuzzF(f *testing.F) {
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FuzzF"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FuzzF"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FuzzF"
		_ = err
	}
}

func TestFunctionLiteral(t *testing.T) {
	testsetup()
	t.Run("test", func(t *testing.T) {
		ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
		_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
		_ = err
		if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in anonymous function"
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
