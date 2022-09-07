//golangcitest:args -Eonlyany
package testdata

import "fmt"

func onlyanyTest() {
	var a interface{} // want `use any instead of an empty interface`
	fmt.Print(a)
}
