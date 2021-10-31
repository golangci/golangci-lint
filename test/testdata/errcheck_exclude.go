//args: -Eerrcheck
//config: linters-settings.errcheck.check-blank=true
//config: linters-settings.errcheck.exclude=testdata/errcheck/exclude.txt
package testdata

import (
	"io/ioutil"
)

func TestErrcheckExclude() []byte {
	ret, _ := ioutil.ReadFile("f.txt")
	return ret
}

func TestErrcheckNoExclude() []byte {
	ret, _ := ioutil.ReadAll(nil) // ERROR "Error return value of `ioutil.ReadAll` is not checked"
	return ret
}
