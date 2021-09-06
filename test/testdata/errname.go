//args: -Eerrname
package testdata

import (
	"errors"
	"fmt"
	"strings"
)

var (
	EOF          = errors.New("end of file")
	ErrEndOfFile = errors.New("end of file")
	errEndOfFile = errors.New("end of file")

	EndOfFileError = errors.New("end of file") // ERROR "the variable name `EndOfFileError` should conform to the `ErrXxx` format"
	ErrorEndOfFile = errors.New("end of file") // ERROR "the variable name `ErrorEndOfFile` should conform to the `ErrXxx` format"
	EndOfFileErr   = errors.New("end of file") // ERROR "the variable name `EndOfFileErr` should conform to the `ErrXxx` format"
	endOfFileError = errors.New("end of file") // ERROR "the variable name `endOfFileError` should conform to the `errXxx` format"
	errorEndOfFile = errors.New("end of file") // ERROR "the variable name `errorEndOfFile` should conform to the `errXxx` format"
)

const maxSize = 256

var (
	ErrOutOfSize = fmt.Errorf("out of size (max %d)", maxSize)
	errOutOfSize = fmt.Errorf("out of size (max %d)", maxSize)

	OutOfSizeError = fmt.Errorf("out of size (max %d)", maxSize) // ERROR "the variable name `OutOfSizeError` should conform to the `ErrXxx` format"
	outOfSizeError = fmt.Errorf("out of size (max %d)", maxSize) // ERROR "the variable name `outOfSizeError` should conform to the `errXxx` format"
)

func errInsideFuncIsNotSentinel() error {
	var lastErr error
	return lastErr
}

type NotErrorType struct{}

func (t NotErrorType) Set() {}
func (t NotErrorType) Get() {}

type DNSConfigError struct{}

func (D DNSConfigError) Error() string { return "DNS config error" }

type someTypeWithoutPtr struct{}           // ERROR "the type name `someTypeWithoutPtr` should conform to the `xxxError` format"
func (s someTypeWithoutPtr) Error() string { return "someTypeWithoutPtr" }

type SomeTypeWithoutPtr struct{}           // ERROR "the type name `SomeTypeWithoutPtr` should conform to the `XxxError` format"
func (s SomeTypeWithoutPtr) Error() string { return "SomeTypeWithoutPtr" }

type someTypeWithPtr struct{}            // ERROR "the type name `someTypeWithPtr` should conform to the `xxxError` format"
func (s *someTypeWithPtr) Error() string { return "someTypeWithPtr" }

type SomeTypeWithPtr struct{}            // ERROR "the type name `SomeTypeWithPtr` should conform to the `XxxError` format"
func (s *SomeTypeWithPtr) Error() string { return "SomeTypeWithPtr" }

type ValidationErrors []string

func (ve ValidationErrors) Error() string { return strings.Join(ve, "\n") }

type validationErrors []string

func (ve validationErrors) Error() string { return strings.Join(ve, "\n") }

type TenErrors [10]string

func (te TenErrors) Error() string { return strings.Join(te[:], "\n") }

type tenErrors [10]string

func (te tenErrors) Error() string { return strings.Join(te[:], "\n") }

type MultiError []error             // ERROR "the type name `MultiError` should conform to the `XxxErrors` format"
func (me MultiError) Error() string { return "" }

type multiError []error             // ERROR "the type name `multiError` should conform to the `xxxErrors` format"
func (me multiError) Error() string { return "" }

type TwoError [2]error            // ERROR "the type name `TwoError` should conform to the `XxxErrors` format"
func (te TwoError) Error() string { return te[0].Error() + "\n" + te[1].Error() }

type twoError [2]error            // ERROR "the type name `twoError` should conform to the `xxxErrors` format"
func (te twoError) Error() string { return te[0].Error() + "\n" + te[1].Error() }
