//golangcitest:args -Etypecheck
package testdata

func TypeCheckBadCalls() {
	typecheckNotExists1.F1() // want "undeclared name: `typecheckNotExists1`"
	typecheckNotExists2.F2() // want "undeclared name: `typecheckNotExists2`"
	typecheckNotExists3.F3() // want "undeclared name: `typecheckNotExists3`"
	typecheckNotExists4.F4() // want "undeclared name: `typecheckNotExists4`"
}
