//args: -Edepguard
//config_path: testdata/configs/depguard.yml
package testdata

import (
	"compress/gzip" // ERROR "`compress/gzip` is in the blacklist"
	"log"           // ERROR "`log` is in the blacklist: don't use log"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
