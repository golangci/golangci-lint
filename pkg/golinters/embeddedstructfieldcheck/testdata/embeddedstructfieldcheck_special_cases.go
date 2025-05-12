//golangcitest:args -Eembeddedstructfieldcheck
package simple

import "time"

func myFunction() {
	type myType struct {
		version   int
		time.Time // want `embedded fields should be listed before regular fields`
	}
}
