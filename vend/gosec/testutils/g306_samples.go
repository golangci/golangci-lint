package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG306 - Poor permissions for WriteFile
var SampleCodeG306 = []CodeSample{
	{[]string{`
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/dat1", d1, 0744)
	check(err)

	allowed := ioutil.WriteFile("/tmp/dat1", d1, 0600)
	check(allowed)

	f, err := os.Create("/tmp/dat2")
	check(err)

	defer f.Close()

	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)

	defer check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	n3, err := f.WriteString("writes\n")
	fmt.Printf("wrote %d bytes\n", n3)

	f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", n4)

	w.Flush()

}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	content := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/dat1", content, os.ModePerm)
	check(err)
}
`}, 1, gosec.NewConfig()},
}
