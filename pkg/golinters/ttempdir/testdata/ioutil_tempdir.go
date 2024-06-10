//golangcitest:args -Ettempdir
package testdata

import (
	"io/ioutil"
	"testing"
)

var (
	_, ee = ioutil.TempDir("a", "b") // never seen
)

func setup() {
	ioutil.TempDir("a", "b")           // never seen
	_, err := ioutil.TempDir("a", "b") // never seen
	if err != nil {
		_ = err
	}
	ioutil.TempDir("a", "b") // never seen
}

func F(t *testing.T) {
	setup()
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in F"
		_ = err
	}
}

func BF(b *testing.B) {
	TBF(b)
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `b\\.TempDir\\(\\)` in BF"
		_ = err
	}
}

func TBF(tb testing.TB) {
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `tb\\.TempDir\\(\\)` in TBF"
		_ = err
	}
}

func FF(f *testing.F) {
	ioutil.TempDir("a", "b")           // want "ioutil\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_, err := ioutil.TempDir("a", "b") // want "ioutil\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
	_ = err
	if _, err := ioutil.TempDir("a", "b"); err != nil { // want "ioutil\\.TempDir\\(\\) should be replaced by `f\\.TempDir\\(\\)` in FF"
		_ = err
	}
}
