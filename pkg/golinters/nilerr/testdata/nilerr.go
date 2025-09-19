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

func l() error {
	aChan := make(chan error, 1)
	bChan := make(chan error, 1)

	var aErr error
	var bErr error

	for i := 0; i < 2; i++ {
		select {
		case err := <-aChan:
			aErr = err
		case err := <-bChan:
			bErr = err
		}
	}

	if aErr != nil {
		return nil // want `error is not nil \(lines \[41 45\]\) but it returns nil`
	}
	if bErr != nil {
		return nil // want `error is not nil \(lines \[42 45\]\) but it returns nil`
	}

	return nil
}
