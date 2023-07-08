//golangcitest:args -Eunexportedglobal
package testdata

import "fmt"

var ExportedVar = 1 // ok

var unexportedVar = 1 // want `unexported global "unexportedVar" should be prefixed with '_'`

var _unexportedVar = 1 // ok

func _1() {
	var local = 42
	fmt.Println(local) // ok
}
