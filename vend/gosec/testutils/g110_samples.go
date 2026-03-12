package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG110 - potential DoS vulnerability via decompression bomb
var SampleCodeG110 = []CodeSample{
	{[]string{`
package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
)

func main() {
	buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
		47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
	b := bytes.NewReader(buff)

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		panic(err)
	}

	r.Close()
}`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
)

func main() {
	buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
		47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
	b := bytes.NewReader(buff)

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 8)
	_, err = io.CopyBuffer(os.Stdout, r, buf)
	if err != nil {
		panic(err)
	}
	r.Close()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"archive/zip"
	"io"
	"os"
	"strconv"
)

func main() {
	r, err := zip.OpenReader("tmp.zip")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for i, f := range r.File {
		out, err := os.OpenFile("output" + strconv.Itoa(i), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		rc, err := f.Open()
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(out, rc)

		out.Close()
		rc.Close()

		if err != nil {
			panic(err)
		}
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"io"
	"os"
)

func main() {
	s, err := os.Open("src")
	if err != nil {
		panic(err)
	}
	defer s.Close()

	d, err := os.Create("dst")
	if err != nil {
		panic(err)
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if  err != nil {
		panic(err)
	}
}
`}, 0, gosec.NewConfig()},
}
