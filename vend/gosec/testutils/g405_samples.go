package testutils

import "github.com/securego/gosec/v2"

var (
	// SampleCodeG405 - Use of weak crypto encryption DES
	SampleCodeG405 = []CodeSample{
		{[]string{`
package main

import (
	"crypto/des"
	"fmt"
)

func main() {
	// Weakness: Usage of weak encryption algorithm

	c, e := des.NewCipher([]byte("mySecret"))

	if e != nil {
		panic("We have a problem: " + e.Error())
	}

	data := []byte("hello world")
	fmt.Println("Plain", string(data))
	c.Encrypt(data, data)

	fmt.Println("Encrypted", string(data))
	c.Decrypt(data, data)

	fmt.Println("Plain Decrypted", string(data))
}

`}, 1, gosec.NewConfig()},
	}

	// SampleCodeG405b - Use of weak crypto encryption RC4
	SampleCodeG405b = []CodeSample{
		{[]string{`
package main

import (
	"crypto/rc4"
	"fmt"
)

func main() {
	// Weakness: Usage of weak encryption algorithm
	
	c, _ := rc4.NewCipher([]byte("mySecret"))

	data := []byte("hello world")
	fmt.Println("Plain", string(data))
	c.XORKeyStream(data, data)

	cryptCipher2, _ := rc4.NewCipher([]byte("mySecret"))

	fmt.Println("Encrypted", string(data))
	cryptCipher2.XORKeyStream(data, data)

	fmt.Println("Plain Decrypted", string(data))
}

`}, 2, gosec.NewConfig()},
	}
)
