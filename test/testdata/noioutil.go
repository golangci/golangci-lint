// args: -Enoioutil
package testdata

import (
	"io/ioutil" // ERROR "io/ioutil package is used"
)

func noIoUtil() {
	_ = ioutil.Discard
}
