package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG503 - Blocklisted import RC4
var SampleCodeG503 = []CodeSample{
	{[]string{`
package main

import (
	"crypto/rc4"
	"encoding/hex"
	"fmt"
)

func main() {
	cipher, err := rc4.NewCipher([]byte("sekritz"))
	if err != nil {
		panic(err)
	}
	plaintext := []byte("I CAN HAZ SEKRIT MSG PLZ")
	ciphertext := make([]byte, len(plaintext))
	cipher.XORKeyStream(ciphertext, plaintext)
	fmt.Println("Secret message is: %s", hex.EncodeToString(ciphertext))
}
`}, 1, gosec.NewConfig()},
}
