//golangcitest:args -Estaticcheck
//golangcitest:config_path testdata/gosimple.yml
package testdata

import (
	"log"
)

func Gosimple(ss []string) {
	if ss != nil { // want "S1031: unnecessary nil check around range"
		for _, s := range ss {
			log.Printf(s)
		}
	}
}

func GosimpleNolintGosimple(ss []string) {
	if ss != nil { //nolint:staticcheck
		for _, s := range ss {
			log.Printf(s)
		}
	}
}
