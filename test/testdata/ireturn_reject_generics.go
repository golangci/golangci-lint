//golangcitest:args -Eireturn
//golangcitest:config_path testdata/configs/ireturn_reject_generics.yml
package testdata

import (
	"bytes"
	"io"
)

func NewWriter() io.Writer {
	var buf bytes.Buffer
	return &buf
}

func TestError() error {
	return nil
}

func Get[K comparable, V int64 | float64](m map[K]V) V { // want `Get returns generic interface \(V\)`
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
