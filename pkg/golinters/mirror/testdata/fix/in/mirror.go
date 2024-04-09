//golangcitest:args -Emirror
//golangcitest:expected_exitcode 0
package testdata

import (
	"unicode/utf8"
)

func foobar() {
	_ = utf8.RuneCount([]byte("foobar"))
}
