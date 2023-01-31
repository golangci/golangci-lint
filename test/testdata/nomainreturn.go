//golangcitest:args -Enomainreturn
package testdata

func main() {
	myFunc := func() {
		ok := true
		if !ok {
			return
		}
	}
	myFunc()

	if e := err(); e != nil {
		return // want "return found in main"
	} else if true {
		return // want "return found in main"
	} else if false {
		return // want "return found in main"
	} else {
		return // want "return found in main"
	}

	return // want "return found in main"
}

func err() error {
	return nil
}
