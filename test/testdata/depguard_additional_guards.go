//golangcitest:args -Edepguard
//golangcitest:config_path testdata/configs/depguard_additional_guards.yml
package testdata

import (
	"compress/gzip" // ERROR "`compress/gzip` is in the denylist"
	"fmt"           // ERROR "`fmt` is in the denylist"
	"log"           // ERROR "`log` is in the denylist: don't use log"
	"strings"       // ERROR "`strings` is in the denylist: disallowed in additional guard"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
	log.Println(fmt.Sprintf("SpewDebugInfo"))
	log.Println(strings.ToLower("SpewDebugInfo"))
}
