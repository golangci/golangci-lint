//args: -Evarerr
package testdata

func f() {
	var err error // ERROR "the error type `err` is initialized with var"
	print(err.Error())

	if err := func() error {
		var err error // novarerr
		return err
	}(); err != nil {
		print(err)
	}
}
