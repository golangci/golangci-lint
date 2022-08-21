//golangcitest:args -Eerrorlint
package testdata

import (
	"errors"
	"fmt"
	"log"
)

var errLintFoo = errors.New("foo")

type errLintBar struct{}

func (*errLintBar) Error() string {
	return "bar"
}

func errorLintAll() {
	err := func() error { return nil }()
	if err == errLintFoo { // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		log.Println("errCompare")
	}

	err = errors.New("oops")
	fmt.Errorf("error: %v", err) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"

	switch err.(type) { // want "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *errLintBar:
		log.Println("errLintBar")
	}
}
