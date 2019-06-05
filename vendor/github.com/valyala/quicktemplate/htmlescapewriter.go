package quicktemplate

import (
	"bytes"
	"io"
)

type htmlEscapeWriter struct {
	w io.Writer
}

func (w *htmlEscapeWriter) Write(b []byte) (int, error) {
	if bytes.IndexByte(b, '<') < 0 &&
		bytes.IndexByte(b, '>') < 0 &&
		bytes.IndexByte(b, '"') < 0 &&
		bytes.IndexByte(b, '\'') < 0 &&
		bytes.IndexByte(b, '&') < 0 {

		// fast path - nothing to escape
		return w.w.Write(b)
	}

	// slow path
	write := w.w.Write
	j := 0
	for i, c := range b {
		switch c {
		case '<':
			write(b[j:i])
			write(strLT)
			j = i + 1
		case '>':
			write(b[j:i])
			write(strGT)
			j = i + 1
		case '"':
			write(b[j:i])
			write(strQuot)
			j = i + 1
		case '\'':
			write(b[j:i])
			write(strApos)
			j = i + 1
		case '&':
			write(b[j:i])
			write(strAmp)
			j = i + 1
		}
	}
	if n, err := write(b[j:]); err != nil {
		return j + n, err
	}
	return len(b), nil
}

var (
	strLT   = []byte("&lt;")
	strGT   = []byte("&gt;")
	strQuot = []byte("&quot;")
	strApos = []byte("&#39;")
	strAmp  = []byte("&amp;")
)
