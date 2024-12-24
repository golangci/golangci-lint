//golangcitest:args -Egoimports
//golangcitest:config_path testdata/goimports_local.yml
package testdata

import (
	"fmt"

	"github.com/golangci/golangci-lint/pkg/config" // want "File is not properly formatted"
	"golang.org/x/tools/go/analysis"
)

func GoimportsLocalPrefixTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = analysis.Analyzer{}
}
