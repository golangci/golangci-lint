//args: -Egci
//config: linters-settings.gci.local-prefixes=github.com/golangci/golangci-lint
package gci

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
