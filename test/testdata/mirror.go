//golangcitest:args -Emirror
package testdata

import (
	"strings"
	"unicode/utf8"
)

func foobar() {
	_ = utf8.RuneCount([]byte("foobar"))                                                                            // want `avoid allocations with utf8\.RuneCountInString`
	_ = strings.Compare(string([]byte{'f', 'o', 'o', 'b', 'a', 'r'}), string([]byte{'f', 'o', 'o', 'b', 'a', 'r'})) // want `avoid allocations with bytes\.Compare`
}
