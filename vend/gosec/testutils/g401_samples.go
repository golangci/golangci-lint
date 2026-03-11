package testutils

import "github.com/securego/gosec/v2"

var (
	// SampleCodeG401 - Use of weak crypto hash MD5
	SampleCodeG401 = []CodeSample{
		{[]string{`
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	defer func() {
	  err := f.Close()
	  if err != nil {
		 log.Printf("error closing the file: %s", err)
	  }
	}()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x", h.Sum(nil))
}
`}, 1, gosec.NewConfig()},
	}

	// SampleCodeG401b - Use of weak crypto hash SHA1
	SampleCodeG401b = []CodeSample{
		{[]string{`
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
)
func main() {
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x", h.Sum(nil))
}
`}, 1, gosec.NewConfig()},
	}
)
