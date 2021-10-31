//args: -Eerrcheck
//config: linters-settings.errcheck.check-blank=true
//config: linters-settings.errcheck.exclude=testdata/errcheck/exclude.txt
package testdata

import (
	"io"
	"os"
)

func TestErrcheckExclude() []byte {
	ret, _ := os.ReadFile("f.txt")
	return ret
}

func TestErrcheckNoExclude() []byte {
	ret, _ := io.ReadAll(nil) // ERROR "Error return value of `io.ReadAll` is not checked"
	return ret
}
