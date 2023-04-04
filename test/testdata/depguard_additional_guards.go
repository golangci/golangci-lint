//golangcitest:args -Edepguard
//golangcitest:config_path testdata/configs/depguard_additional_guards.yml
package testdata

import (
	"compress/gzip" // want "`compress/gzip` is in the denylist"
	"fmt"           // want "`fmt` is in the denylist"
	"log"           // want "`log` is in the denylist: don't use log"
	"strings"       // want "`strings` is in the denylist: disallowed in additional guard"
)

func SpewDebugInfo() {
	log.Println(gzip.BestCompression)
	log.Println(fmt.Sprintf("SpewDebugInfo"))
	log.Println(strings.ToLower("SpewDebugInfo"))
}
