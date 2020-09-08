//args: -Egovet
//config: linters-settings.govet.enable=["ifaceassert"]
package govet

import (
	"io"
)

func IfaceAssert() {
	var v interface {
		Read()
	}
	_ = v.(io.Reader) // ERROR "composites: `os.PathError` composite literal uses unkeyed fields"
}
