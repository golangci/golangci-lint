//golangcitest:args -Eifshort --internal-cmd-test
package testdata

func DontUseShortSyntaxWhenPossible() {
	getValue := func() interface{} { return nil }

	v := getValue() // want "variable 'v' is only used in the if-statement .*"
	if v != nil {
		return
	}
}
