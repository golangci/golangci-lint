package testdata

import (
	"net/http"
	"os"
)

func RetErr() error {
	return nil
}

func MissedErrorCheck() {
	RetErr() // ERROR "Error return value of `RetErr` is not checked"
}

func IgnoreCloseMissingErrHandling() error {
	f, err := os.Open("t.go")
	if err != nil {
		return err
	}

	f.Close()
	return nil
}

func IgnoreCloseInDeferMissingErrHandling() {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	panic(resp)
}

func IgnoreStdxWrite() {
	os.Stdout.Write([]byte{})
	os.Stderr.Write([]byte{})
}
