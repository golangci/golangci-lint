//golangcitest:args -Egci
//golangcitest:config_path testdata/configs/gci.yml
package testdata

import (
	"fmt"

	"github.com/anduril/golangci-lint/pkg/config" // ERROR "File is not \\`gci\\`-ed with --skip-generated -s standard,prefix\\(github.com/anduril/golangci-lint\\),default"

	"github.com/pkg/errors" // ERROR "File is not \\`gci\\`-ed with --skip-generated -s standard,prefix\\(github.com/anduril/golangci-lint\\),default"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = errors.New("")
}
