//golangcitest:args -Edepguard
//golangcitest:config_path testdata/configs/depguard_ignore_file_rules.yml
//golangcitest:expected_exitcode 0
package testdata

// NOTE - No lint errors because this file is ignored
import (
	"compress/gzip"
	"log"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
