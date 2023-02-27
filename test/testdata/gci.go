//golangcitest:args -Egci
//golangcitest:config_path testdata/configs/gci.yml
package testdata

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/config" // want "File is not \\`gci\\`-ed with --skip-generated -s standard,prefix\\(github.com/golangci/golangci-lint\\),default"

	"golang.org/x/tools/go/analysis" // want "File is not \\`gci\\`-ed with --skip-generated -s standard,prefix\\(github.com/golangci/golangci-lint\\),default"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = analysis.Analyzer{}
}
