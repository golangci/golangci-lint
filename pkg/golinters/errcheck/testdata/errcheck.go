//golangcitest:args -Eerrcheck
package testdata

import (
	"bytes"
	"net/http"
	"os"
)

func RetErr() error {
	return nil
}

func MissedErrorCheck() {
	RetErr() // want "Error return value is not checked"
}

func IgnoreCloseMissingErrHandling() error {
	f, err := os.Open("t.go")
	if err != nil {
		return err
	}

	f.Close() // want "Error return value of `f.Close` is not checked"
	return nil
}

func IgnoreCloseInDeferMissingErrHandling() {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // want "Error return value of `resp.Body.Close` is not checked"

	panic(resp)
}

func IgnoreStdxWrite() {
	os.Stdout.Write([]byte{}) // want "Error return value of `os.Stdout.Write` is not checked"
	os.Stderr.Write([]byte{}) // want "Error return value of `os.Stderr.Write` is not checked"
}

func IgnoreBufferWrites(buf *bytes.Buffer) {
	buf.WriteString("x")
}
