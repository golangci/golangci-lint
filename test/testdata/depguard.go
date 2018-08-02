// args: -Edepguard --depguard.include-go-root --depguard.packages='compress/*,log'
package testdata

import (
	"compress/gzip" // ERROR "`compress/gzip` is in the blacklist"
	"log"           // ERROR "`log` is in the blacklist"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
}
