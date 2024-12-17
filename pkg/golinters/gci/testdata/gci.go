//golangcitest:args -Egci
//golangcitest:config_path testdata/gci.yml
package testdata

// want +1 "Invalid import order"
import (
	"golang.org/x/tools/go/analysis"
	"github.com/golangci/golangci-lint/pkg/config"
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
