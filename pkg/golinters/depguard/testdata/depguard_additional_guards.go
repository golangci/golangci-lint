//golangcitest:args -Edepguard
//golangcitest:config_path testdata/depguard_additional_guards.yml
package testdata

import (
	"compress/gzip" // want "import 'compress/gzip' is not allowed from list 'main': nope"
	"fmt"           // want "import 'fmt' is not allowed from list 'main': nope"
	"log"           // want "import 'log' is not allowed from list 'main': don't use log"
	"strings"       // want "import 'strings' is not allowed from list 'main': nope"

	"golang.org/x/tools/go/analysis" // want "import 'golang.org/x/tools/go/analysis' is not allowed from list 'main': example import with dot"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
	log.Println(fmt.Sprintf("SpewDebugInfo"))
	log.Println(strings.ToLower("SpewDebugInfo"))
	_ = analysis.Analyzer{}
}
