//golangcitest:args -Egovet
//golangcitest:config_path testdata/govet.yml
package testdata

import (
	"fmt"
	"io"
	"os"
)

func GovetComposites() error {
	return &os.PathError{"first", "path", os.ErrNotExist} // want "composites: io/fs\\.PathError struct literal uses unkeyed fields"
}

func GovetShadow(f io.Reader, buf []byte) (err error) {
	if f != nil {
		_, err := f.Read(buf) // want `shadow: declaration of .err. shadows declaration at line \d+`
		if err != nil {
			return err
		}
	}
	// Use variable to trigger shadowing error
	_ = err
	return
}

func GovetNolintVet() error {
	return &os.PathError{"first", "path", os.ErrNotExist} //nolint:vet
}

func GovetNolintVetShadow() error {
	return &os.PathError{"first", "path", os.ErrNotExist} //nolint:vetshadow
}

func GovetPrintf() {
	x := "dummy"
	fmt.Printf("%d", x) // want "printf: fmt.Printf format %d has arg x of wrong type string"
}

func GovetStringIntConv() {
	i := 42
	fmt.Println("i = " + string(i)) // want "stringintconv: conversion from int to string yields a string of one rune, not a string of digits \\(did you mean fmt.Sprint\\(x\\)\\?\\)"
}
