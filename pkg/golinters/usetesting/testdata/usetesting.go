//golangcitest:args -Eusetesting
package testdata

import (
	"os"
	"testing"
)

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "") // want `os\.MkdirTemp\(\) could be replaced by <t/b/tb>\.TempDir\(\) in .+`
}

func Test_osSetenv(t *testing.T) {
	os.Setenv("", "")
}
