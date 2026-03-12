package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG403 - weak key strength
var SampleCodeG403 = []CodeSample{
	{[]string{`
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func main() {
	//Generate Private Key
	pvk, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pvk)
}
`}, 1, gosec.NewConfig()},
}
