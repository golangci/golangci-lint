package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestCodeClimate_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "minor",
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
			Severity:   "major",
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

	printer := NewCodeClimate(log, buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `[{"description":"linter-a: some issue","check_name":"linter-a","severity":"minor","fingerprint":"BA73C5DF4A6FD8462FFF1D3140235777","location":{"path":"path/to/filea.go","lines":{"begin":10}}},{"description":"linter-b: another issue","check_name":"linter-b","severity":"major","fingerprint":"0777B4FE60242BD8B2E9B7E92C4B9521","location":{"path":"path/to/fileb.go","lines":{"begin":300}}},{"description":"linter-c: without severity","check_name":"linter-c","severity":"critical","fingerprint":"84F89700554F16CCEB6C0BB481B44AD2","location":{"path":"path/to/filec.go","lines":{"begin":300}}},{"description":"linter-d: unknown severity","check_name":"linter-d","severity":"critical","fingerprint":"340BB02E7B583B9727D73765CB64F56F","location":{"path":"path/to/filed.go","lines":{"begin":300}}}]
`

	assert.Equal(t, expected, buf.String())
}
