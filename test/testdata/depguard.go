//args: -Edepguard
//config: linters-settings.depguard.include-go-root=true
//config: linters-settings.depguard.packages=compress/*,log
package testdata

import (
	"compress/gzip" // ERROR "`compress/gzip` is in the blacklist"
	"log"           // ERROR "`log` is in the blacklist"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
