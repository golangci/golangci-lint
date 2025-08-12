//golangcitest:args -Eusetesting
//golangcitest:config_path testdata/usetesting_configuration.yml
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
	os.Chdir("")
}

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "")
}

func Test_osTempDir(t *testing.T) {
	os.TempDir() // want `os\.TempDir\(\) could be replaced by t\.TempDir\(\) in .+`
}

func Test_osSetenv(t *testing.T) {
	os.Setenv("", "") // want `os\.Setenv\(\) could be replaced by t\.Setenv\(\) in .+`
}

func Test_osCreateTemp(t *testing.T) {
	os.CreateTemp("", "")
	os.CreateTemp("", "xx")
	os.CreateTemp(os.TempDir(), "xx") // want `os\.TempDir\(\) could be replaced by t\.TempDir\(\) in .+`
	os.CreateTemp(t.TempDir(), "xx")
}
