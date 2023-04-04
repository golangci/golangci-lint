//golangcitest:args -Enilerr
package testdata

import "os"

func nilErr1() error {
	err := nilErrDo()
	if err == nil {
		return err // want `error is nil \(line 7\) but it returns error`
	}

	return nil
}

func nilErr2() error {
	err := nilErrDo()
	if err == nil {
		return err // want `error is nil \(line 16\) but it returns error`
	}

	return nil
}

func nilErr3() error {
	err := nilErrDo()
	if err != nil {
		return nil // want `error is not nil \(line 25\) but it returns nil`
	}

	return nil
}

func nilErrDo() error {
	return os.ErrNotExist
}
