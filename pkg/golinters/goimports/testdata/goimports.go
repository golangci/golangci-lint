//golangcitest:config_path testdata/goimports.yml
package testdata

import (
	"fmt" // want "File is not properly formatted"
	"github.com/golangci/golangci-lint/v2/pkg/config"
)

func Bar() {
	fmt.Print("x")
	_ = config.Config{}
}
