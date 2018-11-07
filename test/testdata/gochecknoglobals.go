//args: -Egochecknoglobals
package testdata

import (
	"errors"
	"fmt"
)

var noGlobalsVar int // ERROR "`noGlobalsVar` is a global variable"
var ErrSomeType = errors.New("test that global erorrs aren't warned")

func NoGlobals() {
	fmt.Print(noGlobalsVar)
}
