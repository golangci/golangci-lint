package p

import (
	"net/http"
	"os"
)

func retErr() error {
	return nil
}

func missedErrorCheck() {
	retErr() // ERROR "Error return value is not checked"
}

func ignoreCloseMissingErrHandling() error {
	f, err := os.Open("t.go")
	if err != nil {
		return err
	}

	f.Close()
	return nil
}

func ignoreCloseInDeferMissingErrHandling() {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	panic(resp)
}
