//golangcitest:args -Egoimports
//golangcitest:config_path testdata/configs/goimports.yml
package goimports

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/pkg/errors"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = errors.New("")
}
