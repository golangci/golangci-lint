//golangcitest:args -Egci
//golangcitest:config_path testdata/configs/gci.yml
//golangcitest:expected_exitcode 0
package gci

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/golangci/golangci-lint/pkg/config"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = errors.New("")
}
