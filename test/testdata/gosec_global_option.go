//golangcitest:args -Egosec
//golangcitest:config_path testdata/configs/gosec_global_option.yml
package testdata

import (
	"crypto/md5" // want "G501: Blocklisted import crypto/md5: weak cryptographic primitive"
	"log"
)

func Gosec() {
	// #nosec G401
	h := md5.New() // want "G401: Use of weak cryptographic primitive"
	log.Print(h)
}
