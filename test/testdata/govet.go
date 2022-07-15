//golangcitest:args -Egovet
//golangcitest:config linters-settings.govet.check-shadowing=true
package testdata

import (
	"fmt"
	"io"
	"os"
)

func GovetComposites() error {
	return &os.PathError{"first", "path", os.ErrNotExist} // ERROR "composites: \\`(os|io/fs)\\.PathError\\` composite literal uses unkeyed fields"
}

func GovetShadow(f io.Reader, buf []byte) (err error) {
	if f != nil {
		_, err := f.Read(buf) // ERROR `shadow: declaration of .err. shadows declaration at line \d+`
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
	fmt.Printf("%d", x) // ERROR "printf: fmt.Printf format %d has arg x of wrong type string"
}

func GovetStringIntConv() {
	i := 42
	fmt.Println("i = " + string(i)) // ERROR "stringintconv: conversion from int to string yields a string of one rune, not a string of digits \\(did you mean fmt.Sprint\\(x\\)\\?\\)"
}
