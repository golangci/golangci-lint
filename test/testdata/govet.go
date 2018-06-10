// args: -Egovet --govet.check-shadowing=true
package testdata

import (
	"io"
	"os"
)

func Govet() error {
	return &os.PathError{"first", "path", os.ErrNotExist} // ERROR "os.PathError composite literal uses unkeyed fields"
}

func GovetShadow(f io.Reader, buf []byte) (err error) {
	if f != nil {
		_, err := f.Read(buf) // ERROR "declaration of .err. shadows declaration at testdata/govet.go:\d+"
		if err != nil {
			return err
		}
	}
	// Use variable to trigger shadowing error
	_ = err
	return
}
