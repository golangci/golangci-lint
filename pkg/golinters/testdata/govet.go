package testdata

import "os"

func govet() error {
	return &os.PathError{"first", "path", os.ErrNotExist} // ERROR "os.PathError composite literal uses unkeyed fields"
}

func govetShadow(f *os.File, buf []byte) (err error) {
	if f != nil {
		_, err := f.Read(buf) // ERROR "declaration of .err. shadows declaration at testdata/govet.go:9"
		if err != nil {
			return err
		}
	}
	// Use variable to trigger shadowing error
	_ = err
	return
}
