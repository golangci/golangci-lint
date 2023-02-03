//golangcitest:args -Eunusedtypeparam
package testdata

type Constraint interface {
	string | ~int
}

func ok[E Constraint](arg E) {
	arg2 := arg
	_ = arg2
	var arg3 E
	_ = arg3
}

func ng[E Constraint](arg any) { // want "This func unused type parameter."
	arg2 := arg
	_ = arg2
}
