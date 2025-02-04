//golangcitest:args -Eusetesting
package testdata

import (
	"os"
	"testing"
)

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "") // want `os\.MkdirTemp\(\) could be replaced by t\.TempDir\(\) in .+`
}

func Test_osSetenv(t *testing.T) {
	os.Setenv("", "") // want `os\.Setenv\(\) could be replaced by t\.Setenv\(\) in .+`
}

func Test_osTempDir(t *testing.T) {
	os.TempDir()
}

func Test_osCreateTemp(t *testing.T) {
	os.CreateTemp("", "")   // want `os\.CreateTemp\("", \.\.\.\) could be replaced by os\.CreateTemp\(t\.TempDir\(\), \.\.\.\) in .+`
	os.CreateTemp("", "xx") // want `os\.CreateTemp\("", \.\.\.\) could be replaced by os\.CreateTemp\(t\.TempDir\(\), \.\.\.\) in .+`
	os.CreateTemp(os.TempDir(), "xx")
	os.CreateTemp(t.TempDir(), "xx")
}
