package testutils

import "github.com/securego/gosec/v2"

var (
	// SampleCodeG406 - Use of deprecated weak crypto hash MD4
	SampleCodeG406 = []CodeSample{
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

	// SampleCodeG406b - Use of deprecated weak crypto hash RIPEMD160
	SampleCodeG406b = []CodeSample{
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
)
