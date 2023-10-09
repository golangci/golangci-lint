//golangcitest:args -Eperfsprint
package testdata

import "fmt"

func SprintfCouldBeStrconv() {
	fmt.Sprintf("%d", 42) // want "Sprintf can be replaced with faster strconv.Itoa"
}
