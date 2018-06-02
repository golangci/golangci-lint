package withtests

import "fmt"

var varUsedOnlyInTests bool

func usedOnlyInTests() {}

type someType struct {
	fieldUsedOnlyInTests bool
	fieldUsedHere        bool
}

func usedHere() {
	v := someType{
		fieldUsedHere: true,
	}
	fmt.Println(v)
}

func init() {
	usedHere()
}
