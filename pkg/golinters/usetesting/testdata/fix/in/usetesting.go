//golangcitest:args -Eusetesting
//golangcitest:expected_exitcode 0
package testdata

import (
	"os"
	"testing"
)

func Test_osCreateTemp(t *testing.T) {
	os.CreateTemp("", "")   // want `os\.CreateTemp\("", \.\.\.\) could be replaced by os\.CreateTemp\(t\.TempDir\(\), \.\.\.\) in .+`
	os.CreateTemp("", "xx") // want `os\.CreateTemp\("", \.\.\.\) could be replaced by os\.CreateTemp\(t\.TempDir\(\), \.\.\.\) in .+`
	os.CreateTemp(os.TempDir(), "xx")
	os.CreateTemp(t.TempDir(), "xx")
}
