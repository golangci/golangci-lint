//golangcitest:args -Egci
//golangcitest:config_path testdata/configs/gci.yml
//golangcitest:expected_exitcode 0
package gci

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/pkg/errors"
	"fmt"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = errors.New("")
}
