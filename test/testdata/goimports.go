//golangcitest:args -Egoimports
package testdata

import (
	"fmt" // want "File is not `goimports`-ed"
	"github.com/golangci/golangci-lint/pkg/config"
)

func Bar() {
	fmt.Print("x")
	_ = config.Config{}
}
