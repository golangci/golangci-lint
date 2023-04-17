//golangcitest:args -Enakedefer
//golangcitest:config_path testdata/configs/nakedefer.yml
package testdata

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func funcNotReturnAnyType() {
}

func funcReturnErr() error {
	return errors.New("some error")
}

func funcReturnFuncAndErr() (func(), error) {
	return func() {
	}, nil
}

func ignoreFunc() error {
	return errors.New("some error")
}

func testCaseValid1() {
	defer funcNotReturnAnyType() // ignore

	defer func() { // ignore
		funcNotReturnAnyType()
	}()

	defer func() { // ignore
		_ = funcReturnErr()
	}()
}

func testCaseInvalid1() {
	defer funcReturnErr() // want "deferred call should not return anything"

	defer funcReturnFuncAndErr() // want "deferred call should not return anything"

	defer func() error { // want "deferred call should not return anything"
		return nil
	}()

	defer func() func() { // want "deferred call should not return anything"
		return func() {}
	}()
}

func testCase1() {
	defer fmt.Errorf("some text") // want "deferred call should not return anything"

	r := new(bytes.Buffer)
	defer io.LimitReader(r, 1) // want "deferred call should not return anything"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("DONE"))
	}))
	defer srv.Close()                  // ignore
	defer srv.CloseClientConnections() // ignore
	defer srv.Certificate()            // want "deferred call should not return anything"
}

func testCaseExclude1() {
	// exclude ignoreFunc
	defer ignoreFunc() // ignore
}

func testCaseExclude2() {
	// exclude os\.(Create|WriteFile|Chmod)
	defer os.Create("file_test1")                                   // ignore
	defer os.WriteFile("file_test2", []byte("data"), os.ModeAppend) // ignore
	defer os.Chmod("file_test3", os.ModeAppend)                     // ignore
	defer os.FindProcess(100500)                                    // want "deferred call should not return anything"
}

func testCaseExclude3() {
	// exclude fmt\.Print.*
	defer fmt.Println("e1")        // ignore
	defer fmt.Print("e1")          // ignore
	defer fmt.Printf("e1")         // ignore
	defer fmt.Sprintf("some text") // want "deferred call should not return anything"
}

func testCaseExclude4() {
	// exclude io\.Close
	rc, _ := zlib.NewReader(bytes.NewReader([]byte("111")))
	defer rc.Close() // ignore
}
