//go:build go1.20

//golangcitest:args -Etypecheck
package testdata

func TypeCheckBadCalls() {
	typecheckNotExists1.F1() // want "undefined: typecheckNotExists1"
	typecheckNotExists2.F2() // want "undefined: typecheckNotExists2"
	typecheckNotExists3.F3() // want "undefined: typecheckNotExists3"
	typecheckNotExists4.F4() // want "undefined: typecheckNotExists4"
}
