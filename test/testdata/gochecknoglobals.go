//golangcitest:args -Egochecknoglobals
package testdata

import (
	"errors"
	"fmt"
	"regexp"
)

var noGlobalsVar int // ERROR "noGlobalsVar is a global variable"
var ErrSomeType = errors.New("test that global erorrs aren't warned")

var (
	OnlyDigites = regexp.MustCompile(`^\d+$`)
	BadNamedErr = errors.New("this is bad") // ERROR "BadNamedErr is a global variable"
)

func NoGlobals() {
	fmt.Print(noGlobalsVar)
}
