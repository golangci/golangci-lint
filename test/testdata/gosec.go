//args: -Egosec
package testdata

import (
	"crypto/md5" // ERROR "G501: Blacklisted import `crypto/md5`: weak cryptographic primitive"
	"log"
)

func Gosec() {
	h := md5.New() // ERROR "G401: Use of weak cryptographic primitive"
	log.Print(h)
}

func GosecNolintGas() {
	h := md5.New() //nolint:gas
	log.Print(h)
}

func GosecNolintGosec() {
	h := md5.New() //nolint:gosec
	log.Print(h)
}
