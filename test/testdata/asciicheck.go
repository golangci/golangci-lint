//golangcitest:args -Easciicheck
package testdata

import (
	"fmt"
	"time"
)

type AsciicheckTеstStruct struct { // ERROR `identifier "AsciicheckTеstStruct" contain non-ASCII character: U\+0435 'е'`
	Date time.Time
}

type AsciicheckField struct{}

type AsciicheckJustStruct struct {
	Tеst AsciicheckField // ERROR `identifier "Tеst" contain non-ASCII character: U\+0435 'е'`
}

func AsciicheckTеstFunc() { // ERROR `identifier "AsciicheckTеstFunc" contain non-ASCII character: U\+0435 'е'`
	var tеstVar int // ERROR `identifier "tеstVar" contain non-ASCII character: U\+0435 'е'`
	tеstVar = 0
	fmt.Println(tеstVar)
}
