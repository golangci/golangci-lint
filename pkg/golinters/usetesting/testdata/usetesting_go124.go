//go:build go1.24

//golangcitest:args -Eusetesting
package testdata

import (
	"context"
	"os"
	"testing"
)

func Test_contextBackground(t *testing.T) {
	context.Background() // want `context\.Background\(\) could be replaced by <t/b/tb>\.Context\(\) in .+`
}

func Test_contextTODO(t *testing.T) {
	context.TODO() // want `context\.TODO\(\) could be replaced by <t/b/tb>\.Context\(\) in .+`
}

func Test_osChdir(t *testing.T) {
	os.Chdir("") // want `os\.Chdir\(\) could be replaced by <t/b/tb>\.Chdir\(\) in .+`
}

func Test_osMkdirTemp(t *testing.T) {
	os.MkdirTemp("", "") // want `os\.MkdirTemp\(\) could be replaced by <t/b/tb>\.TempDir\(\) in .+`
}

func Test_osSetenv(t *testing.T) {
	os.Setenv("", "")
}
