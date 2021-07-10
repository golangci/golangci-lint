//args: -Eerrcheck
//config: linters-settings.errcheck.check-blank=true
//config: linters-settings.errcheck.exclude-functions=io/ioutil.ReadFile,io/ioutil.ReadDir
package testdata

import (
	"io/ioutil"
)

func TestErrcheckExclude() []byte {
	ret, _ := ioutil.ReadFile("f.txt")
	_, _ = ioutil.ReadDir("dir")
	return ret
}

func TestErrcheckNoExclude() []byte {
	ret, _ := ioutil.ReadAll(nil) // ERROR "Error return value of `ioutil.ReadAll` is not checked"
	return ret
}
