package testdata

import (
	"crypto/md5" // ERROR "G501: Blacklisted import crypto/md5: weak cryptographic primitive"
	"log"        // nolint:depguard
)

func Gas() {
	h := md5.New() // ERROR "G401: Use of weak cryptographic primitive"
	log.Print(h)
}
