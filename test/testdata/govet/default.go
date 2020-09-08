//args: -Egovet
//config: linters-settings.govet.check-shadowing=true
package govet

import (
	"fmt"
	"io"
	"os"
)

func Composites() error {
	return &os.PathError{"first", "path", os.ErrNotExist} // ERROR "composites: \\`(os|io/fs)\\.PathError\\` composite literal uses unkeyed fields"
}

func Shadow(f io.Reader, buf []byte) (err error) {
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

func NolintVet() error {
	return &os.PathError{"first", "path", os.ErrNotExist} //nolint:vet
}

func NolintVetShadow() error {
	return &os.PathError{"first", "path", os.ErrNotExist} //nolint:vetshadow
}

func Printf() {
	x := "dummy"
	fmt.Printf("%d", x) // ERROR "printf: Printf format %d has arg x of wrong type string"
}

func StringIntConv() {
	i := 42
	fmt.Println("i = " + string(i)) // ERROR "stringintconv: conversion from int to string yields a string of one rune, not a string of digits (did you mean fmt.Sprint(x)?)"
}
