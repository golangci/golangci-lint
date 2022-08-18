//golangcitest:args -Eerrname
package testdata

import (
	"errors"
	"fmt"
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
