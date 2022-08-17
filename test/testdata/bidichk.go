//golangcitest:args -Ebidichk
package testdata

import "fmt"

func main() {
	fmt.Println("LEFT-TO-RIGHT-OVERRIDE: '‭', it is between the single quotes, but it is not visible with a regular editor") // ERROR "found dangerous unicode character sequence LEFT-TO-RIGHT-OVERRIDE"
}
