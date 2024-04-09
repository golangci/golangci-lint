//golangcitest:args -Egci
//golangcitest:config_path testdata/gci.yml
package testdata

import (
	"golang.org/x/tools/go/analysis" // want "File is not \\`gci\\`-ed with --skip-generated -s standard -s prefix\\(github.com/golangci/golangci-lint,github.com/daixiang0/gci\\) -s default --custom-order"
	"github.com/golangci/golangci-lint/pkg/config"
	"fmt"
	"errors"
	gcicfg "github.com/daixiang0/gci/pkg/config"  // want "File is not \\`gci\\`-ed with --skip-generated -s standard -s prefix\\(github.com/golangci/golangci-lint,github.com/daixiang0/gci\\) -s default --custom-order"
)

func GoimportsLocalTest() {
	fmt.Print(errors.New("x"))
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
