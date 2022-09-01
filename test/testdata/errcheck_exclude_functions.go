//golangcitest:args -Eerrcheck
//golangcitest:config_path testdata/configs/exclude_functions.yml
package testdata

import (
	"io/ioutil"
)

func TestErrcheckExcludeFunctions() []byte {
	ret, _ := ioutil.ReadFile("f.txt")
	ioutil.ReadDir("dir")
	return ret
}

func TestErrcheckNoExcludeFunctions() []byte {
	ret, _ := ioutil.ReadAll(nil) // want "Error return value of `ioutil.ReadAll` is not checked"
	return ret
}
