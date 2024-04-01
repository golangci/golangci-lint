//golangcitest:args -Edepguard
//golangcitest:config_path testdata/depguard.yml
package testdata

import (
	"compress/gzip" // want "import 'compress/gzip' is not allowed from list 'main': nope"
	"log"           // want "import 'log' is not allowed from list 'main': don't use log"

	"golang.org/x/tools/go/analysis" // want "import 'golang.org/x/tools/go/analysis' is not allowed from list 'main': example import with dot"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
	_ = analysis.Analyzer{}
}
