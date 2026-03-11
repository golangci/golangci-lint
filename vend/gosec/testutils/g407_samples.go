package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG407 - Use of hardcoded nonce/IV
var SampleCodeG407 = []CodeSample{
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesOFB := cipher.NewOFB(block, []byte("ILoveMyNonceAlot"))
	var output = make([]byte, 16)
	aesOFB.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func encrypt(nonce []byte) {
	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesOFB := cipher.NewOFB(block, nonce)
	var output = make([]byte, 16)
	aesOFB.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))
}

func main() {

	var nonce = []byte("ILoveMyNonceAlot")
	encrypt(nonce)
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesOFB := cipher.NewOFB(block, []byte("ILoveMyNonceAlot")) // #nosec G407
	var output = make([]byte, 16)
	aesOFB.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

}

`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher( []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)

	cipherText := aesGCM.Seal(nil, func() []byte {
		if true {
			return []byte("ILoveMyNonce")
		} else {
			return []byte("IDont'Love..")
		}
	}(), []byte("My secret message"), nil) // #nosec G407
	fmt.Println(string(cipherText))

	cipherText, _ = aesGCM.Open(nil, func() []byte {
		if true {
			return []byte("ILoveMyNonce")
		} else {
			return []byte("IDont'Love..")
		}
	}(), cipherText, nil) // #nosec G407

	fmt.Println(string(cipherText))
}
`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesOFB := cipher.NewOFB(block, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	var output = make([]byte, 16)
	aesOFB.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

}`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCTR := cipher.NewCTR(block, []byte("ILoveMyNonceAlot"))
	var output = make([]byte, 16)
	aesCTR.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

}`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCTR := cipher.NewCTR(block, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	var output = make([]byte, 16)
	aesCTR.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)

	cipherText := aesGCM.Seal(nil, []byte("ILoveMyNonce"), []byte("My secret message"), nil)
	fmt.Println(string(cipherText))
	cipherText, _ = aesGCM.Open(nil, []byte("ILoveMyNonce"), cipherText, nil)
	fmt.Println(string(cipherText))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)
	cipherText := aesGCM.Seal(nil, []byte{}, []byte("My secret message"), nil)
	fmt.Println(string(cipherText))

	cipherText, _ = aesGCM.Open(nil, []byte{}, cipherText, nil)
	fmt.Println(string(cipherText))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)

	cipherText := aesGCM.Seal(nil, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []byte("My secret message"), nil)
	fmt.Println(string(cipherText))

	cipherText, _ = aesGCM.Open(nil, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, cipherText, nil)
	fmt.Println(string(cipherText))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher( []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)

	cipherText := aesGCM.Seal(nil, func() []byte {
		if true {
			return []byte("ILoveMyNonce")
		} else {
			return []byte("IDont'Love..")
		}
	}(), []byte("My secret message"), nil)
	fmt.Println(string(cipherText))

	cipherText, _ = aesGCM.Open(nil, func() []byte {
		if true {
			return []byte("ILoveMyNonce")
		} else {
			return []byte("IDont'Love..")
		}
	}(), cipherText, nil)

	fmt.Println(string(cipherText))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)

	cipherText := aesGCM.Seal(nil, func() []byte {
		if true {
			return []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		} else {
			return []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		}
	}(), []byte("My secret message"), nil)
	fmt.Println(string(cipherText))

	cipherText, _ = aesGCM.Open(nil, func() []byte {
		if true {
			return []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		} else {
			return []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		}
	}(), cipherText, nil)
	fmt.Println(string(cipherText))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)
	cipheredText := aesGCM.Seal(nil, func() []byte { return []byte("ILoveMyNonce") }(), []byte("My secret message"), nil)
	fmt.Println(string(cipheredText))
	cipheredText, _ = aesGCM.Open(nil, func() []byte { return []byte("ILoveMyNonce") }(), cipheredText, nil)
	fmt.Println(string(cipheredText))

}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)
	cipheredText := aesGCM.Seal(nil, func() []byte { return []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} }(), []byte("My secret message"), nil)
	fmt.Println(string(cipheredText))
	cipheredText, _ = aesGCM.Open(nil, func() []byte { return []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} }(), cipheredText, nil)
	fmt.Println(string(cipheredText))

}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCFB := cipher.NewCFBEncrypter(block, []byte("ILoveMyNonceAlot"))
	var output = make([]byte, 16)
	aesCFB.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))
	aesCFB = cipher.NewCFBDecrypter(block, []byte("ILoveMyNonceAlot"))
	aesCFB.XORKeyStream(output, output)
	fmt.Println(string(output))

}`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCFB := cipher.NewCFBEncrypter(block, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	var output = make([]byte, 16)
	aesCFB.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))
	aesCFB = cipher.NewCFBDecrypter(block, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCFB.XORKeyStream(output, output)
	fmt.Println(string(output))

}`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCBC := cipher.NewCBCEncrypter(block, []byte("ILoveMyNonceAlot"))

	var output = make([]byte, 16)
	aesCBC.CryptBlocks(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

	aesCBC = cipher.NewCBCDecrypter(block, []byte("ILoveMyNonceAlot"))
	aesCBC.CryptBlocks(output, output)
	fmt.Println(string(output))

}`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCBC := cipher.NewCBCEncrypter(block, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})

	var output = make([]byte, 16)
	aesCBC.CryptBlocks(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))

	aesCBC = cipher.NewCBCDecrypter(block, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCBC.CryptBlocks(output, output)
	fmt.Println(string(output))

}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	var nonce = []byte("ILoveMyNonce")
	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)
	fmt.Println(string(aesGCM.Seal(nil, nonce, []byte("My secret message"), nil)))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func main() {

	var nonce = []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCTR := cipher.NewCTR(block, nonce)
	var output = make([]byte, 16)
	aesCTR.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func coolFunc(size int) []byte{
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}

func main() {

	var nonce = coolFunc(16)
	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesCTR := cipher.NewCTR(block, nonce)
	var output = make([]byte, 16)
	aesCTR.XORKeyStream(output, []byte("Very Cool thing!"))
	fmt.Println(string(output))
}
`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var nonce = []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

func main() {

	block, _ := aes.NewCipher([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	aesGCM, _ := cipher.NewGCM(block)
	cipherText := aesGCM.Seal(nil, nonce, []byte("My secret message"), nil)
	fmt.Println(string(cipherText))

}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, nil
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {}
`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

const iv = "1234567812345678"

func wrapper(s string, b cipher.Block) {
	cipher.NewCTR(b, []byte(s))
}

func main() {
	b, _ := aes.NewCipher([]byte("1234567812345678"))
	wrapper(iv, b)
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

var globalIV = []byte("1234567812345678")

func wrapper(iv []byte, b cipher.Block) {
	cipher.NewCTR(b, iv)
}

func main() {
	b, _ := aes.NewCipher([]byte("1234567812345678"))
	wrapper(globalIV, b)
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/cipher"
)

func recursive(s string, b cipher.Block) {
	recursive(s, b)
	cipher.NewCTR(b, []byte(s))
}

func main() {
	recursive("1234567812345678", nil)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

func main() {
	k := make([]byte, 48)
	key, iv := k[:32], k[32:]
	block, _ := aes.NewCipher(key)
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

func main() {
	k := make([]byte, 48)
	k[32] = 1
	key, iv := k[:32], k[32:]
	block, _ := aes.NewCipher(key)
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv)
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func main() {
	iv := make([]byte, 16)
	io.ReadFull(nil, iv)
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

func fill(b []byte) {
	b[0] = 1
}

func main() {
	iv := make([]byte, 16)
	fill(iv)
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv)
	iv[0] = 1 // overwriting
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv[0:8])
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv[0:16])
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	buf := make([]byte, 128)
	rand.Read(buf[32:48])
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, buf[32:48])
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
)

func main() {
	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	iv := []byte("1234567890123456")

	var f func(cipher.Block, []byte) cipher.Stream
	if len(os.Args) > 1 {
		f = cipher.NewCTR
	} else {
		f = cipher.NewOFB
	}
	stream := f(block, iv)
	_ = stream
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"os"
)

func main() {
	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	iv := []byte("1234567890123456")
	rand.Read(iv)

	var f func(cipher.Block, []byte) cipher.Stream
	if len(os.Args) > 1 {
		f = cipher.NewCTR
	} else {
		f = cipher.NewOFB
	}
	stream := f(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (

"crypto/aes"
"crypto/cipher"
"crypto/rand"

)

func myReaderDirect(b []byte) (int, error) {
	return rand.Read(b)
}

func main() {
	iv := make([]byte, 16)
	// Direct call to user function (myReaderDirect) which calls rand.Read
	myReaderDirect(iv)

	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	_ = cipher.NewCTR(block, iv)
}
	   `}, 0, gosec.NewConfig()},
	{[]string{`package main

import (

"crypto/aes"
"crypto/cipher"
"crypto/rand"

)

func myReaderDirect(b []byte) (int, error) {
	n, err := rand.Read(b)
	if n > 1 {
		b[0] = 1 // overwriting
	}
	return n, err
}

func main() {
	iv := make([]byte, 16)
	// Direct call to user function (myReaderDirect) which calls rand.Read but overwrites the IV
	myReaderDirect(iv)

	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	_ = cipher.NewCTR(block, iv)
}
	   `}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
"crypto/cipher"
)

func myBadCipher(n int, block cipher.Block) cipher.Stream {
    iv := make([]byte, n) 
    iv[0] = 0x01
    return cipher.NewCTR(block, iv)
}
	   `}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
"crypto/cipher"
)

func myBadCipher(n int, block cipher.Block) cipher.Stream {
    iv := make([]byte, n) 
    return cipher.NewCTR(block, iv)
}
	   `}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
"crypto/cipher"
"os"
)

func myGoodCipher(block cipher.Block) (cipher.Stream, error) {
    iv, err := os.ReadFile("iv.bin")
    if err != nil {
        return nil, err
    }
    return cipher.NewCTR(block, iv), nil
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
"crypto/cipher"
"io"
)

func myGoodInterfaceCipher(r io.Reader, block cipher.Block) {
    iv := make([]byte, 16)
    r.Read(iv)
    stream := cipher.NewCTR(block, iv)
    _ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	iv := []byte("1234567890123456")
	iv[8] = 0
	rand.Read(iv)
	stream := cipher.NewCTR(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
)

func test(init func([]byte)) {
	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	iv := make([]byte, 16)
	init(iv) // We can't resolve 'init', should default to Dynamic to avoid FP
	stream := cipher.NewCTR(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},

	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type CustomReader interface {
	io.Reader
}

func testCustomReader(cr CustomReader) {
	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	iv := make([]byte, 16)
	cr.Read(iv)
	stream := cipher.NewCTR(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/cipher"
	"io"
)

func interfaceSafeOverwrite(r io.Reader, block cipher.Block) {
	iv := make([]byte, 16)
	iv[0] = 0 // Tainted
	r.Read(iv) // Dynamic Interface Read (covers taint)
	stream := cipher.NewCTR(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type CustomReader interface {
	io.Reader
}

func testCustomReaderOverwrite(cr CustomReader) {
	key := []byte("example key 1234")
	block, _ := aes.NewCipher(key)
	iv := make([]byte, 16)
	iv[15] = 1 // Taint
	cr.Read(iv) // Cover via embedded interface
	stream := cipher.NewCTR(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/cipher"
	"io"
)

func interfaceSafeOverwriteSlice(r io.Reader, block cipher.Block) {
	iv := make([]byte, 16)
	iv[0] = 0
	r.Read(iv[:])
	stream := cipher.NewCTR(block, iv)
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/cipher"
)

func pointerUnOpIV(block cipher.Block) {
	iv := make([]byte, 16) // Hardcoded
	ptr := &iv
	stream := cipher.NewCTR(block, *ptr)
	_ = stream
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/rand"
	"crypto/cipher"
)

func pointerUnOpSafeIV(block cipher.Block) {
	iv := make([]byte, 16)
	rand.Read(iv) // Dynamic
	ptr := &iv
	stream := cipher.NewCTR(block, *ptr) // dynamic dereference
	_ = stream
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv[6:12])
	rand.Read(iv[0:6])
	rand.Read(iv[12:16])
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv[6:12])
	iv[6] = 0
	rand.Read(iv[0:7])
	iv[10] = 0
	rand.Read(iv[10:16])
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
    "os"
)
func main() {
    iv := make([]byte, len(os.Args))
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	// Slice with Variable Bound (Unresolvable Range)
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "os"
)
func main() {
    iv := make([]byte, 16)
    low := len(os.Args)
    sub := iv[low:]
    rand.Read(sub)
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	// IndexAddr with Variable Index (Unresolvable Range)
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
    "os"
)
func main() {
    iv := make([]byte, 16)
    i := len(os.Args)
    iv[i] = 0
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
)
func test(iv []byte) {
    iv[0] = 0
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
func main() {
    test(make([]byte, 16))
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv[6:12])
	iv[6] = 0
	rand.Read(iv[0:7])
	iv[10] = 0
	rand.Read(iv[10:16])
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
    "os"
)
func main() {
    iv := make([]byte, len(os.Args))
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "os"
)
func main() {
    iv := make([]byte, 16)
    low := len(os.Args)
    sub := iv[low:]
    rand.Read(sub)
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
    "os"
)
func main() {
    iv := make([]byte, 16)
    i := len(os.Args)
    iv[i] = 0
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main
import (
    "crypto/aes"
    "crypto/cipher"
)
func test(iv []byte) {
    iv[0] = 0
    block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
    _ = cipher.NewCTR(block, iv)
}
func main() {
    test(make([]byte, 16))
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func unsafeOverwrite(i int) {
	iv := make([]byte, 16)
	rand.Read(iv)
	if i >= 10 && i < 16 {
		iv[i] = 0
	}
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv[:16])
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func safeOverwrite(i int) {
	iv := make([]byte, 128)
	rand.Read(iv)
	if i >= 16 && i < 128{
		iv[i] = 0
	}
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv[:16])
}
`}, 0, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func unsafeOverwrite(i int) {
	iv := make([]byte, 16)
	rand.Read(iv)
	if i > 0 {
		iv[i % 16] = 0
	}
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func unsafeOverwrite(i int) {
	iv := make([]byte, 16)
	rand.Read(iv)
	if i - 16 > 0 && i + 16 < 32 {
		iv[i - 16] = 0
	}
	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
	{[]string{`package main
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func main() {
	iv := make([]byte, 16)
	rand.Read(iv)
	
	// Alias assignment
	alias := iv
	alias[0] = 0 // Hardcoded write via alias (unsafe)

	block, _ := aes.NewCipher([]byte("12345678123456781234567812345678"))
	_ = cipher.NewCTR(block, iv)
}
`}, 1, gosec.NewConfig()},
}
