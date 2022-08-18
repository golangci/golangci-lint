//golangcitest:args -Egovet
//golangcitest:config_path testdata/configs/govet_ifaceassert.yml
package testdata

import (
	"io"
)

func GovetIfaceAssert() {
	var v interface {
		Read()
	}
	_ = v.(io.Reader) // ERROR "impossible type assertion: no type can implement both interface\\{Read\\(\\)\\} and io\\.Reader \\(conflicting types for Read method\\)"
}
