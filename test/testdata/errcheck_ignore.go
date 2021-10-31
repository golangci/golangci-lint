//args: -Eerrcheck
//config_path: testdata/errcheck/ignore_config.yml
package testdata

import (
	"fmt"
	"io"
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
	ret, _ := os.ReadFile("f.txt")
	return ret
}

func TestErrcheckNoIgnoreIoutil() []byte {
	ret, _ := io.ReadAll(nil) // ERROR "Error return value of `io.ReadAll` is not checked"
	return ret
}
