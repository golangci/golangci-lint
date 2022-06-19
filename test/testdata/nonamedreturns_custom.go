// args: -Enonamedreturns
// config_path: testdata/configs/nonamedreturns.yml
package testdata

func simple() (err error) {
	defer func() {
		err = nil
	}()
	return
}

func twoReturnParams() (i int, err error) { // ERROR `named return "i" with type "int" found`
	defer func() {
		i = 0
		err = nil
	}()
	return
}

func allUnderscoresExceptError() (_ int, err error) {
	defer func() {
		err = nil
	}()
	return
}

func customName() (myName error) {
	defer func() {
		myName = nil
	}()
	return
}

func errorIsNoAssigned() (err error) { // ERROR `named return "err" with type "error" found`
	defer func() {
		_ = err
		processError(err)
		if err == nil {
		}
		switch err {
		case nil:
		default:
		}
	}()
	return
}

func shadowVariable() (err error) {
	defer func() {
		err := 123 // linter doesn't understand that this is different variable (even if different type) (yet?)
		_ = err
	}()
	return
}

func shadowVariable2() (err error) {
	defer func() {
		a, err := doSomething() // linter doesn't understand that this is different variable (yet?)
		_ = a
		_ = err
	}()
	return
}

type myError = error // linter doesn't understand that this is the same type (yet?)

func customType() (err myError) { // ERROR `named return "err" with type "myError" found`
	defer func() {
		err = nil
	}()
	return
}

func notTheLast() (err error, _ int) {
	defer func() {
		err = nil
	}()
	return
}

func twoErrorsCombined() (err1, err2 error) {
	defer func() {
		err1 = nil
		err2 = nil
	}()
	return
}

func twoErrorsSeparated() (err1 error, err2 error) {
	defer func() {
		err1 = nil
		err2 = nil
	}()
	return
}

func errorSlice() (err []error) { // ERROR `named return "err" with type "\[\]error" found`
	defer func() {
		err = nil
	}()
	return
}

func deferWithVariable() (err error) { // ERROR `named return "err" with type "error" found`
	f := func() {
		err = nil
	}
	defer f() // linter can't catch closure passed via variable (yet?)
	return
}

func uberMultierr() (err error) { // ERROR `named return "err" with type "error" found`
	defer func() {
		multierrAppendInto(&err, nil) // linter doesn't allow it (yet?)
	}()
	return
}

func deferInDefer() (err error) {
	defer func() {
		defer func() {
			err = nil
		}()
	}()
	return
}

func twoDefers() (err error) {
	defer func() {}()
	defer func() {
		err = nil
	}()
	return
}

func callFunction() (err error) {
	defer func() {
		_, err = doSomething()
	}()
	return
}

func callFunction2() (err error) {
	defer func() {
		var a int
		a, err = doSomething()
		_ = a
	}()
	return
}

func deepInside() (err error) {
	if true {
		switch true {
		case false:
			for i := 0; i < 10; i++ {
				go func() {
					select {
					default:
						defer func() {
							if true {
								switch true {
								case false:
									for j := 0; j < 10; j++ {
										go func() {
											select {
											default:
												err = nil
											}
										}()
									}
								}
							}
						}()
					}
				}()
			}
		}
	}
	return
}

var goodFuncLiteral = func() (err error) {
	defer func() {
		err = nil
	}()
	return
}

var badFuncLiteral = func() (err error) { // ERROR `named return "err" with type "error" found`
	defer func() {
		_ = err
	}()
	return
}

func funcLiteralInsideFunc() error {
	do := func() (err error) {
		defer func() {
			err = nil
		}()
		return
	}
	return do()
}

type x struct{}

func (x) goodMethod() (err error) {
	defer func() {
		err = nil
	}()
	return
}

func (x) badMethod() (err error) { // ERROR `named return "err" with type "error" found`
	defer func() {
		_ = err
	}()
	return
}

func processError(error)                    {}
func doSomething() (int, error)             { return 10, nil }
func multierrAppendInto(*error, error) bool { return false } // https://pkg.go.dev/go.uber.org/multierr#AppendInto
