//args: -Eerrcheck
//config_path: testdata/errcheck/exclude_functions.yml
package testdata

import (
	"io"
	"os"
)

func TestErrcheckExcludeFunctions() []byte {
	ret, _ := os.ReadFile("f.txt")
	os.ReadDir("dir")
	return ret
}

func TestErrcheckNoExcludeFunctions() []byte {
	ret, _ := io.ReadAll(nil) // ERROR "Error return value of `io.ReadAll` is not checked"
	return ret
}
