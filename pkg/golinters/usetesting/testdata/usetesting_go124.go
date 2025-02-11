//go:build go1.24

//golangcitest:args -Eusetesting
package testdata

import (
	"context"
	"os"
	"testing"
)

func Test_contextBackground(t *testing.T) {
	context.Background() // want `context\.Background\(\) could be replaced by t\.Context\(\) in .+`
}

func Test_contextTODO(t *testing.T) {
	context.TODO() // want `context\.TODO\(\) could be replaced by t\.Context\(\) in .+`
}

func Test_osChdir(t *testing.T) {
	os.Chdir("") // want `os\.Chdir\(\) could be replaced by t\.Chdir\(\) in .+`
}

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
