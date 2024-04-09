//golangcitest:args -Edepguard
//golangcitest:config_path testdata/depguard_ignore_file_rules.yml
//golangcitest:expected_exitcode 0
package testdata

// NOTE - No lint errors because this file is ignored
import (
	"compress/gzip"
	"log"

	"golang.org/x/tools/go/analysis"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
	_ = analysis.Analyzer{}
}
