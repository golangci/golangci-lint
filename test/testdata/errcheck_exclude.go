//golangcitest:args -Eerrcheck
//golangcitest:config_path testdata/configs/errcheck_exclude.yml
package testdata

import (
	"io/ioutil"
)

func TestErrcheckExclude() []byte {
	ret, _ := ioutil.ReadFile("f.txt")
	return ret
}

func TestErrcheckNoExclude() []byte {
	ret, _ := ioutil.ReadAll(nil) // want "Error return value of `ioutil.ReadAll` is not checked"
	return ret
}
