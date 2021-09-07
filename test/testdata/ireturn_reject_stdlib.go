// args: -Eireturn
// config_path: testdata/configs/ireturn_stdlib_reject.yml
package testdata

import (
	"bytes"
	"io"
)

func NewAllowDoer() io.Writer {
	var buf bytes.Buffer
	return &buf
}
