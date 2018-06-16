package testdata

import (
	"bytes"
	"io"
)

func noTransitiveExpanding() {
	//nolint:errcheck
	var buf io.Writer = &bytes.Buffer{}
	buf.Write([]byte("123"))
}

func nolintFuncByInlineCommentDoesNotWork() { //nolint:errcheck
	var buf io.Writer = &bytes.Buffer{}
	buf.Write([]byte("123"))
}
