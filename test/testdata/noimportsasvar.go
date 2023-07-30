//golangcitest:args -Enoimportsasvar
package testdata

import (
	"bytes"
	err "errors"
	"flag"
	"fmt"
	. "io"
	"io/fs"
	"math"
	_ "net/http"
	"os"
)

func sample() {
	// prevent imports being auto removed by go fmt
	_ = fs.ErrClosed
	_ = math.E
	_ = bytes.ErrTooLarge
	_ = flag.ErrHelp
	_ = err.New("fake error")
	_ = os.ErrClosed

	const fs = "FS" // want "const name 'fs' shared with import 'io/fs'"

	var (
		err error // want "var name 'err' shared with import 'err'"
	)

	fmt := fmt.Println // want "var name 'fmt' shared with import 'fmt'"

	for _, bytes := range []byte{0, 1, 2} { // want "var name 'bytes' shared with import 'bytes'"
		fmt(bytes)
	}

	for flag := 0; flag < 1; flag++ { // want "var name 'flag' shared with import 'flag'"
		fmt(flag)
	}

	http := 200 // OK - underscore import, should not be flagged
	io := EOF   // OK - dot import, should not be flagged

	// avoid staticcheck warnings
	fmt(http)
	fmt(io)
	fmt(err)
	fmt(fs)
	_ = doMath(func(a int, b int) int {
		return a + b
	})

}

func doMath(math func(int, int) int) int { // want "var name 'math' shared with import 'math'"
	return math(1, 2)
}

func retOs() (os string) { // want "var name 'os' shared with import 'os"
	return
}
