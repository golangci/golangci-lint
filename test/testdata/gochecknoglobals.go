//golangcitest:args -Egochecknoglobals
package testdata

import (
	"errors"
	"fmt"
	"regexp"
)

var noGlobalsVar int // want "noGlobalsVar is a global variable"
var ErrSomeType = errors.New("test that global errors aren't warned")

var (
	OnlyDigites = regexp.MustCompile(`^\d+$`)
	BadNamedErr = errors.New("this is bad") // want "BadNamedErr is a global variable"
)

func NoGlobals() {
	fmt.Print(noGlobalsVar)
}
