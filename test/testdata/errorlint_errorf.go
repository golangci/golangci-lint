//golangcitest:args -Eerrorlint
//golangcitest:config_path testdata/configs/errorlint_errorf.yml
package testdata

import (
	"errors"
	"fmt"
)

type customError struct{}

func (customError) Error() string {
	return "oops"
}

func errorLintErrorf() {
	err := errors.New("oops")
	fmt.Errorf("error: %w", err)
	fmt.Errorf("error: %v", err)         // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	fmt.Errorf("%v %v", err, err)        // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	fmt.Errorf("error: %s", err.Error()) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	customError := customError{}
	fmt.Errorf("error: %s", customError.Error()) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	strErr := "oops"
	fmt.Errorf("%v", strErr)
}
