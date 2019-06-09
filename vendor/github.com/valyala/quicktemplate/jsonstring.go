package quicktemplate

import (
	"io"
	"strings"
)

func writeJSONString(w io.Writer, s string) {
	if len(s) > 24 &&
		strings.IndexByte(s, '"') < 0 &&
		strings.IndexByte(s, '\\') < 0 &&
		strings.IndexByte(s, '\n') < 0 &&
		strings.IndexByte(s, '\r') < 0 &&
		strings.IndexByte(s, '\t') < 0 &&
		strings.IndexByte(s, '\f') < 0 &&
		strings.IndexByte(s, '\b') < 0 &&
		strings.IndexByte(s, '<') < 0 &&
		strings.IndexByte(s, '\'') < 0 &&
		strings.IndexByte(s, 0) < 0 {

		// fast path - nothing to escape
		w.Write(unsafeStrToBytes(s))
		return
	}

	// slow path
	write := w.Write
	b := unsafeStrToBytes(s)
	j := 0
	n := len(b)
	if n > 0 {
		// Hint the compiler to remove bounds checks in the loop below.
		_ = b[n-1]
	}
	for i := 0; i < n; i++ {
		switch b[i] {
		case '"':
			write(b[j:i])
			write(strBackslashQuote)
			j = i + 1
		case '\\':
			write(b[j:i])
			write(strBackslashBackslash)
			j = i + 1
		case '\n':
			write(b[j:i])
			write(strBackslashN)
			j = i + 1
		case '\r':
			write(b[j:i])
			write(strBackslashR)
			j = i + 1
		case '\t':
			write(b[j:i])
			write(strBackslashT)
			j = i + 1
		case '\f':
			write(b[j:i])
			write(strBackslashF)
			j = i + 1
		case '\b':
			write(b[j:i])
			write(strBackslashB)
			j = i + 1
		case '<':
			write(b[j:i])
			write(strBackslashLT)
			j = i + 1
		case '\'':
			write(b[j:i])
			write(strBackslashQ)
			j = i + 1
		case 0:
			write(b[j:i])
			write(strBackslashZero)
			j = i + 1
		}
	}
	write(b[j:])
}

var (
	strBackslashQuote     = []byte(`\"`)
	strBackslashBackslash = []byte(`\\`)
	strBackslashN         = []byte(`\n`)
	strBackslashR         = []byte(`\r`)
	strBackslashT         = []byte(`\t`)
	strBackslashF         = []byte(`\u000c`)
	strBackslashB         = []byte(`\u0008`)
	strBackslashLT        = []byte(`\u003c`)
	strBackslashQ         = []byte(`\u0027`)
	strBackslashZero      = []byte(`\u0000`)
)
