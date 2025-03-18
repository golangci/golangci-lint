package exitcodes

import (
	"fmt"
)

const (
	Success = iota
	IssuesFound
	WarningInTest
	Failure
	Timeout
	NoGoFiles
	NoConfigFileDetected
	ErrorWasLogged
	PackagesLoadingFailure
)

type ExitError struct {
	Inner   error
	Message string
	Code    int
}

func (e ExitError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s, %s", e.Message, e.Inner.Error())
	}
	return e.Message
}

func (e ExitError) Unwrap() error { return e.Inner }

var (
	ErrNoGoFiles = &ExitError{
		Message: "no go files to analyze",
		Code:    NoGoFiles,
	}
	ErrFailure = &ExitError{
		Message: "failed to analyze",
		Code:    Failure,
	}
)
