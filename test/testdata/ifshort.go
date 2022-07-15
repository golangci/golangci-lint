//golangcitest:args -Eifshort
package testdata

func DontUseShortSyntaxWhenPossible() {
	getValue := func() interface{} { return nil }

	v := getValue() // ERROR "variable 'v' is only used in the if-statement .*"
	if v != nil {
		return
	}
}
