//go:build go1.24

//golangcitest:args -Eusetesting
//golangcitest:config_path testdata/usetesting_go124_configuration.yml
package testdata

import (
	"context"
	"os"
	"testing"
)

func Test_contextBackground(t *testing.T) {
	context.Background()
}

func Test_contextTODO(t *testing.T) {
	context.TODO()
}

func Test_osChdir(t *testing.T) {
	os.Chdir("")
}

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "")
}

func Test_osSetenv(t *testing.T) {
	os.Setenv("", "") // want `os\.Setenv\(\) could be replaced by t\.Setenv\(\) in .+`
}

func Test_osTempDir(t *testing.T) {
	os.TempDir() // want `os\.TempDir\(\) could be replaced by t\.TempDir\(\) in .+`
}

func Test_osCreateTemp(t *testing.T) {
	os.CreateTemp("", "")
	os.CreateTemp("", "xx")
	os.CreateTemp(os.TempDir(), "xx") // want `os\.TempDir\(\) could be replaced by t\.TempDir\(\) in .+`
	os.CreateTemp(t.TempDir(), "xx")
}
