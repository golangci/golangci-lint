//golangcitest:args -Eusetesting
//golangcitest:config_path testdata/usetesting_configuration.yml
package testdata

import (
	"os"
	"testing"
)

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "")
}

func Test_osSetenv(t *testing.T) {
	os.Setenv("", "") // want `os\.Setenv\(\) could be replaced by <t/b/tb>\.Setenv\(\) in .+`
}
