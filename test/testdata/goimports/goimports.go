//args: -Egoimports
//config: linters-settings.goimports.local-prefixes=github.com/anduril/golangci-lint
package goimports

import (
	"fmt"

	"github.com/anduril/golangci-lint/pkg/config"
	"github.com/pkg/errors"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = errors.New("")
}
