//go:build go1.24

//golangcitest:config_path testdata/gci_go124.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"crypto/sha3"
	"errors"
	"fmt"
)

func _() {
	fmt.Print(errors.New("x"))
	sha3.New224()
}
