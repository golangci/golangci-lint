//golangcitest:args -Eunconvert
package testdata

import "log"

func Unconvert() {
	a := 1
	b := int(a) // want "unnecessary conversion"
	log.Print(b)
}
