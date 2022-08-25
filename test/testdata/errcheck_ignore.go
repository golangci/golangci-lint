//golangcitest:args -Eerrcheck
//golangcitest:config_path testdata/configs/ignore_config.yml
package testdata

import (
	"fmt"
	"io/ioutil"
	"os"
)

func TestErrcheckIgnoreOs() {
	_, _ = os.Open("f.txt")
}

func TestErrcheckIgnoreFmt(s string) int {
	n, _ := fmt.Println(s)
	return n
}

func TestErrcheckIgnoreIoutil() []byte {
	ret, _ := ioutil.ReadFile("f.txt")
	return ret
}

func TestErrcheckNoIgnoreIoutil() []byte {
	ret, _ := ioutil.ReadAll(nil) // want "Error return value of `ioutil.ReadAll` is not checked"
	return ret
}
