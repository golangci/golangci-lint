//golangcitest:args -Egci
//golangcitest:config_path testdata/gci.yml
package testdata

import ( // want "File is not properly formatted"
	"golang.org/x/tools/go/analysis"
	"github.com/golangci/golangci-lint/pkg/config" // want "File is not properly formatted"
	"fmt"
	"errors"
	gcicfg "github.com/daixiang0/gci/pkg/config"
)

func GoimportsLocalTest() {
	fmt.Print(errors.New("x"))
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
