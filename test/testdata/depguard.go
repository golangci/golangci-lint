//golangcitest:args -Edepguard
//golangcitest:config_path testdata/configs/depguard.yml
package testdata

import (
	"compress/gzip" // want "`compress/gzip` is in the denylist"
	"log"           // want "`log` is in the denylist: don't use log"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
