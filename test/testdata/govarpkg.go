package testdata

import (
	"fmt"
	"github.com/alexal/govarpkg/pkg/tst"
)

func test() {
	a := "test"
	fmt.Println(a)
	tst := tst.Name{}
	tst.New("Alex")
}
