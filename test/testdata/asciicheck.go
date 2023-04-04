//golangcitest:args -Easciicheck
package testdata

import (
	"fmt"
	"time"
)

type AsciicheckTеstStruct struct { // want `identifier "AsciicheckTеstStruct" contain non-ASCII character: U\+0435 'е'`
	Date time.Time
}

type AsciicheckField struct{}

type AsciicheckJustStruct struct {
	Tеst AsciicheckField // want `identifier "Tеst" contain non-ASCII character: U\+0435 'е'`
}

func AsciicheckTеstFunc() { // want `identifier "AsciicheckTеstFunc" contain non-ASCII character: U\+0435 'е'`
	var tеstVar int // want `identifier "tеstVar" contain non-ASCII character: U\+0435 'е'`
	tеstVar = 0
	fmt.Println(tеstVar)
}
