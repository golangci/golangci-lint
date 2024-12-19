//golangcitest:args -Enilnesserr
package testdata

func do() error {
	return nil
}

func do2() error {
	return nil
}

func someCall() error {
	err := do()
	if err != nil {
		return err
	}
	err2 := do2()
	if err2 != nil {
		return err // want `return a nil value error after check error`
	}
	return nil
}

func sameCall2() error {
	err := do()
	if err == nil {
		err2 := do2()
		if err2 != nil {
			return err // want `return a nil value error after check error`
		}
		return nil
	}
	return err

}
