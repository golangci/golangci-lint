package internal

import (
	"bytes"
	"testing"

	diffpkg "github.com/sourcegraph/go-diff/diff"
)

func BenchmarkParseDiffLines(b *testing.B) {
	h := &diffpkg.Hunk{
		OrigStartLine: 1,
		Body: []byte(`-old line 1
+new line 1
 unchanged
-old line 2
+new line 2
+added line
`),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		parseDiffLines(h)
	}
}

func BenchmarkHunkChangesParser(b *testing.B) {
	h := &diffpkg.Hunk{
		OrigStartLine: 1,
		Body: []byte(`-deleted line 1
+added line 1
 unchanged
+added line 2
-deleted line 2
+added line 3
`),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		p := hunkChangesParser{}
		p.parse(h)
	}
}

func BenchmarkLargeDiff(b *testing.B) {
	var buf bytes.Buffer
	for i := range 1000 {
		buf.WriteString("-old line ")
		buf.WriteString(string(rune('0' + i%10)))
		buf.WriteString("\n")
		buf.WriteString("+new line ")
		buf.WriteString(string(rune('0' + i%10)))
		buf.WriteString("\n")
	}

	h := &diffpkg.Hunk{
		OrigStartLine: 1,
		Body:          buf.Bytes(),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		p := hunkChangesParser{}
		p.parse(h)
		parseDiffLines(h)
	}
}
