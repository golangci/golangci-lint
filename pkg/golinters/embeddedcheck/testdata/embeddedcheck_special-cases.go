//golangcitest:args -Eembeddedcheck
package simple

import "time"

func myFunction() {
	type myType struct {
		version   int
		time.Time // want `embedded types should be listed before non embedded types`
	}
}
