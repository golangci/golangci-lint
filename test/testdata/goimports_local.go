//golangcitest:args -Egoimports
//golangcitest:config_path testdata/configs/goimports_local.yml
package testdata

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/config" // want "File is not `goimports`-ed with -local github.com/golangci/golangci-lint"
	"github.com/pkg/errors"
)

func GoimportsLocalPrefixTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = errors.New("")
}
