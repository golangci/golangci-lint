//args: -Eerrcheck
//config: linters-settings.errcheck.check-blank=true
package testdata

import (
	"fmt"
	"os"
)

func TestErrcheckIgnoreFmtByDefault(s string) int {
	n, _ := fmt.Println(s)
	return n
}

func TestErrcheckNoIgnoreOs() {
	_, _ = os.Open("f.txt") // ERROR "Error return value of `os.Open` is not checked"
}
