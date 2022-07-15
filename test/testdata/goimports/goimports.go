//golangcitest:args -Egoimports
//golangcitest:config linters-settings.goimports.local-prefixes=github.com/golangci/golangci-lint
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
