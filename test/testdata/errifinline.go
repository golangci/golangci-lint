//golangcitest:args -Eerrifinline
package testdata

func error1() error {
	return nil
}

func error2() (any, error) {
	return nil, nil
}

func errIfInline() {
	err := error1()
	if err != nil { // want `inline err assignment in if initializer`
		_ = err
	}

	if err := error1(); err != nil {
		_ = err
	}

	_, err = error2()
	if err != nil { // want `inline err assignment in if initializer`
		_ = err
	}

	if _, err := error2(); err != nil {
		_ = err
	}

	something, err := error2()
	if err != nil {
		_ = err
	}
	_ = something
}
