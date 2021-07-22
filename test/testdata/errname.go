//args: -Eerrname
package testdata

import (
	"errors"
	"fmt"
)

var (
	EOF          = errors.New("end of file")
	ErrEndOfFile = errors.New("end of file")
	errEndOfFile = errors.New("end of file")

	EndOfFileError = errors.New("end of file") // ERROR "the sentinel error `EndOfFileError` should be of the form ErrXxx"
	ErrorEndOfFile = errors.New("end of file") // ERROR "the sentinel error `ErrorEndOfFile` should be of the form ErrXxx"
	EndOfFileErr   = errors.New("end of file") // ERROR "the sentinel error `EndOfFileErr` should be of the form ErrXxx"
	endOfFileError = errors.New("end of file") // ERROR "the sentinel error `endOfFileError` should be of the form errXxx"
	errorEndOfFile = errors.New("end of file") // ERROR "the sentinel error `errorEndOfFile` should be of the form errXxx"
)

const maxSize = 256

var (
	ErrOutOfSize = fmt.Errorf("out of size (max %d)", maxSize)
	errOutOfSize = fmt.Errorf("out of size (max %d)", maxSize)

	OutOfSizeError = fmt.Errorf("out of size (max %d)", maxSize) // ERROR "the sentinel error `OutOfSizeError` should be of the form ErrXxx"
	outOfSizeError = fmt.Errorf("out of size (max %d)", maxSize) // ERROR "the sentinel error `outOfSizeError` should be of the form errXxx"
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

type someTypeWithoutPtr struct{}           // ERROR "the error type `someTypeWithoutPtr` should be of the form xxxError"
func (s someTypeWithoutPtr) Error() string { return "someTypeWithoutPtr" }

type SomeTypeWithoutPtr struct{}           // ERROR "the error type `SomeTypeWithoutPtr` should be of the form XxxError"
func (s SomeTypeWithoutPtr) Error() string { return "SomeTypeWithoutPtr" }

type someTypeWithPtr struct{}            // ERROR "the error type `someTypeWithPtr` should be of the form xxxError"
func (s *someTypeWithPtr) Error() string { return "someTypeWithPtr" }

type SomeTypeWithPtr struct{}            // ERROR "the error type `SomeTypeWithPtr` should be of the form XxxError"
func (s *SomeTypeWithPtr) Error() string { return "SomeTypeWithPtr" }
