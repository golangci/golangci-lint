package testdata

func TypeCheckBadCalls() {
	typecheckNotExists1.F1() // ERROR "undeclared name: typecheckNotExists1"
	typecheckNotExists2.F2() // ERROR "undeclared name: typecheckNotExists2"
	typecheckNotExists3.F3() // ERROR "undeclared name: typecheckNotExists3"
	typecheckNotExists4.F4()
}
