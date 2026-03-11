package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG506 - Blocklisted import MD4
var SampleCodeG506 = []CodeSample{
	{[]string{`
package main

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/md4"
)

func main() {
	h := md4.New()
	h.Write([]byte("test"))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}
`}, 1, gosec.NewConfig()},
}
