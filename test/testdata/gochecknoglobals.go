//args: -Egochecknoglobals
package testdata

import (
	"errors"
	"fmt"
)

var noGlobalsVar int // ERROR "`noGlobalsVar` is a global variable"
var ErrSomeType = errors.New("test that global erorrs aren't warned")
var ErrFmt1 = fmt.Errorf("test that global errors made with fmt aren't warned")

//var re1 = regexp.MustComplile("/test that regexp aren't warned/")

func NoGlobals() {
	_ = ErrFmt1
	fmt.Print(noGlobalsVar)
}
