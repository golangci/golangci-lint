//args: -Edepguard
//config_path: testdata/configs/depguard_ignore_file_rules.yml
package testdata

// NOTE - No lint errors becuase this file is ignored
import (
	"compress/gzip"
	"log"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
