package withtests

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	v := someType{
		fieldUsedOnlyInTests: true,
	}
	fmt.Println(v, varUsedOnlyInTests)
	usedOnlyInTests()
}
