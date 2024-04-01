//golangcitest:args -Enonamedreturns
package testdata

import "errors"

func simple() (err error) {
	defer func() {
		err = nil
	}()
	return
}

func twoReturnParams() (i int, err error) { // want `named return "i" with type "int" found`
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

func errorIsNoAssigned() (err error) { // want `named return "err" with type "error" found`
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

func shadowVariable() (err error) { // want `named return "err" with type "error" found`
	defer func() {
		err := errors.New("xxx")
		_ = err
	}()
	return
}

func shadowVariableButAssign() (err error) {
	defer func() {
		{
			err := errors.New("xxx")
			_ = err
		}
		err = nil
	}()
	return
}

func shadowVariable2() (err error) { // want `named return "err" with type "error" found`
	defer func() {
		a, err := doSomething()
		_ = a
		_ = err
	}()
	return
}

type errorAlias = error

func errorAliasIsTheSame() (err errorAlias) {
	defer func() {
		err = nil
	}()
	return
}

type myError error // linter doesn't check underlying type (yet?)

func customTypeWithErrorUnderline() (err myError) { // want `named return "err" with type "myError" found`
	defer func() {
		err = nil
	}()
	return
}

type myError2 interface{ error } // linter doesn't check interfaces

func customTypeWithTheSameInterface() (err myError2) { // want `named return "err" with type "myError2" found`
	defer func() {
		err = nil
	}()
	return
}

var _ error = myError3{}

type myError3 struct{} // linter doesn't check interfaces

func (m myError3) Error() string { return "" }

func customTypeImplementingErrorInterface() (err myError3) { // want `named return "err" with type "myError3" found`
	defer func() {
		err = struct{}{}
	}()
	return
}

func shadowErrorType() {
	type error interface { // linter understands that this is not built-in error, even if it has the same name
		Error() string
	}
	do := func() (err error) { // want `named return "err" with type "error" found`
		defer func() {
			err = nil
		}()
		return
	}
	do()
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

func errorSlice() (err []error) { // want `named return "err" with type "\[\]error" found`
	defer func() {
		err = nil
	}()
	return
}

func deferWithVariable() (err error) { // want `named return "err" with type "error" found`
	f := func() {
		err = nil
	}
	defer f() // linter can't catch closure passed via variable (yet?)
	return
}

func uberMultierr() (err error) { // want `named return "err" with type "error" found`
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

var badFuncLiteral = func() (err error) { // want `named return "err" with type "error" found`
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

func (x) badMethod() (err error) { // want `named return "err" with type "error" found`
	defer func() {
		_ = err
	}()
	return
}

func processError(error)                    {}
func doSomething() (int, error)             { return 10, nil }
func multierrAppendInto(*error, error) bool { return false } // https://pkg.go.dev/go.uber.org/multierr#AppendInto
