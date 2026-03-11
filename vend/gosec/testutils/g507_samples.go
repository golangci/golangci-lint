package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG507 - Blocklisted import RIPEMD160
var SampleCodeG507 = []CodeSample{
	{[]string{`
package main

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

func main() {
	h := ripemd160.New()
	h.Write([]byte("test"))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}
`}, 1, gosec.NewConfig()},
}
