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
	os.Setenv("", "") // want `os\.Setenv\(\) could be replaced by <t/b/tb>\.Setenv\(\) in .+`
}
