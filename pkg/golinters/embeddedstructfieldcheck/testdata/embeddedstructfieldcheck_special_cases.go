//golangcitest:args -Eembeddedstructfieldcheck
package testdata

import "time"

func myFunction() {
	type myType struct {
		version   int
		time.Time // want `embedded fields should be listed before regular fields`
	}
}
