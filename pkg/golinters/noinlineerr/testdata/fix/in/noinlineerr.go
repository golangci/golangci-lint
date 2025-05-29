//golangcitest:args -Enoinlineerr
//golangcitest:expected_exitcode 0
package testdata

type MyAliasErr error

type MyCustomError struct{}

func (mc *MyCustomError) Error() string {
	return "error"
}

func doSomething() error {
	return nil
}

func doSmthManyArgs(a, b, c, d int) error {
	return nil
}

func doSmthMultipleReturn() (bool, error) {
	return false, nil
}

func doMyAliasErr() MyAliasErr {
	return nil
}

func doMyCustomErr() *MyCustomError {
	return &MyCustomError{}
}

func invalid() error {
	err := doSomething()
	if err != nil {
		return err
	}

	err = doSmthManyArgs(0,
		0,
		0,
		0,
	)
	if err != nil {
		return err
	}

	err = doMyAliasErr()
	if err != nil {
		return err
	}

	err = doMyCustomErr()
	if err != nil {
		return err
	}

	return nil
}
