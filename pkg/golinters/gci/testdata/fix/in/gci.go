//golangcitest:config_path testdata/gci.yml
//golangcitest:expected_exitcode 0
package gci

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"golang.org/x/tools/go/analysis"

	"fmt"

	gcicfg "github.com/daixiang0/gci/pkg/config"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
