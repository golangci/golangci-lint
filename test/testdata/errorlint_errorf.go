//args: -Eerrorlint
//config_path: testdata/configs/errorlint_errorf.yml
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
	fmt.Errorf("error: %v", err)         // ERROR "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	fmt.Errorf("%v %v", err, err)        // ERROR "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	fmt.Errorf("error: %s", err.Error()) // ERROR "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	customError := customError{}
	fmt.Errorf("error: %s", customError.Error()) // ERROR "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
	strErr := "oops"
	fmt.Errorf("%v", strErr)
}
