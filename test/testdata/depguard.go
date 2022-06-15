//args: -Edepguard
//config_path: testdata/configs/depguard.yml
package testdata

import (
	"compress/gzip" // ERROR "`compress/gzip` is in the denylist"
	"log"           // ERROR "`log` is in the denylist: don't use log"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
