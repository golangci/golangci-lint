package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG505 - Blocklisted import SHA1
var SampleCodeG505 = []CodeSample{
	{[]string{`
package main

import (
	"crypto/sha1"
	"fmt"
	"os"
)

func main() {
	for _, arg := range os.Args {
		fmt.Printf("%x - %s\n", sha1.Sum([]byte(arg)), arg)
	}
}
`}, 1, gosec.NewConfig()},
}
