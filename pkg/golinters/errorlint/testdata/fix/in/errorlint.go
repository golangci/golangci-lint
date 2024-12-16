//golangcitest:args -Eerrorlint
//golangcitest:expected_exitcode 0
package testdata

import (
	"errors"
	"fmt"
)

func Good() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %w", err)
}

func NonWrappingVerb() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %v", err)
}

func NonWrappingTVerb() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %T", err)
}

func DoubleNonWrappingVerb() error {
	err := errors.New("oops")
	return fmt.Errorf("%v %v", err, err)
}

func ErrorOneWrap() error {
	err1 := errors.New("oops1")
	err2 := errors.New("oops2")
	err3 := errors.New("oops3")
	return fmt.Errorf("%v, %w, %v", err1, err2, err3)
}

func ValidNonWrappingTVerb() error {
	err1 := errors.New("oops1")
	err2 := errors.New("oops2")
	err3 := errors.New("oops3")
	return fmt.Errorf("%w, %T, %w", err1, err2, err3)
}

func ErrorMultipleWraps() error {
	err1 := errors.New("oops1")
	err2 := errors.New("oops2")
	err3 := errors.New("oops3")
	return fmt.Errorf("%w, %w, %w", err1, err2, err3)
}

func ErrorMultipleWrapsWithCustomError() error {
	err1 := errors.New("oops1")
	err2 := MyError{}
	err3 := errors.New("oops3")
	return fmt.Errorf("%w, %w, %w", err1, err2, err3)
}

func ErrorStringFormat() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %s", err.Error())
}

func ErrorStringFormatCustomError() error {
	err := MyError{}
	return fmt.Errorf("error: %s", err.Error())
}

func NotAnError() error {
	err := "oops"
	return fmt.Errorf("%v", err)
}

type MyError struct{}

func (MyError) Error() string {
	return "oops"
}

func ErrorIndexReset() error {
	err := errors.New("oops1")
	return fmt.Errorf("%[1]v %d %f %[1]v, %d, %f", err, 1, 2.2)
}

func ErrorIndexResetGood() error {
	err := errors.New("oops1")
	return fmt.Errorf("%[1]w %d %f %[1]w, %d, %f", err, 1, 2.2)
}
