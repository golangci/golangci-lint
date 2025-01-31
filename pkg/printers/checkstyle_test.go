package printers

import (
	"bytes"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestCheckstyle_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "warning",
			Text:       "some issue",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
				Column:   4,
			},
		},
		{
			FromLinter: "linter-b",
			Severity:   "error",
			Text:       "another issue",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/fileb.go",
				Offset:   5,
				Line:     300,
				Column:   9,
			},
		},
		{
			FromLinter: "linter-c",
			Severity:   "",
			Text:       "without severity",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/filec.go",
				Offset:   5,
				Line:     300,
				Column:   10,
			},
		},
		{
			FromLinter: "linter-d",
			Severity:   "foo",
			Text:       "unknown severity",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/filed.go",
				Offset:   5,
				Line:     300,
				Column:   11,
			},
		},
	}

	buf := new(bytes.Buffer)

	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)
	log.SetLevel(logutils.LogLevelDebug)

	printer := NewCheckstyle(log, buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `<?xml version="1.0" encoding="UTF-8"?>

<checkstyle version="5.0">
  <file name="path/to/filea.go">
    <error column="4" line="10" message="some issue" severity="warning" source="linter-a"></error>
  </file>
  <file name="path/to/fileb.go">
    <error column="9" line="300" message="another issue" severity="error" source="linter-b"></error>
  </file>
  <file name="path/to/filec.go">
    <error column="10" line="300" message="without severity" severity="error" source="linter-c"></error>
  </file>
  <file name="path/to/filed.go">
    <error column="11" line="300" message="unknown severity" severity="error" source="linter-d"></error>
  </file>
</checkstyle>
`

	assert.Equal(t, expected, strings.ReplaceAll(buf.String(), "\r", ""))
}
