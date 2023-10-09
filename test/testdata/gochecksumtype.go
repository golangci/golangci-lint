//golangcitest:args -Egochecksumtype
package testdata

import (
	"log"
)

//sumtype:decl
type SumType interface{ isSumType() }

//sumtype:decl
type One struct{} // want "type 'One' is not an interface"

func (One) isSumType() {}

type Two struct{}

func (Two) isSumType() {}

func sumTypeTest() {
	var sum SumType = One{}
	switch sum.(type) { // want "exhaustiveness check failed for sum type.*SumType.*missing cases for Two"
	case One:
	}

	switch sum.(type) { // want "exhaustiveness check failed for sum type.*SumType.*missing cases for Two"
	case One:
	default:
		panic("??")
	}

	log.Println("??")

	switch sum.(type) {
	case One:
	case Two:
	}
}
