//golangcitest:args -Eerrcheck
//golangcitest:config linters-settings.errcheck.check-blank=true
package testdata

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func TestErrcheckIgnoreHashWriteByDefault() []byte {
	h := sha256.New()
	h.Write([]byte("food"))
	return h.Sum(nil)
}

func TestErrcheckIgnoreFmtByDefault(s string) int {
	n, _ := fmt.Println(s)
	return n
}

func TestErrcheckNoIgnoreOs() {
	_, _ = os.Open("f.txt") // ERROR "Error return value of `os.Open` is not checked"
}
