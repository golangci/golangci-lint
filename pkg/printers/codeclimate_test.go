package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestCodeClimate_Print(t *testing.T) {
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
			Text:       "issue c",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"ccc\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/filec.go",
				Offset:   6,
				Line:     200,
				Column:   2,
			},
		},
	}

	buf := new(bytes.Buffer)
	printer := NewCodeClimate(buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	//nolint:lll
	expected := `[{"description":"linter-a: some issue","severity":"warning","fingerprint":"BA73C5DF4A6FD8462FFF1D3140235777","location":{"path":"path/to/filea.go","lines":{"begin":10}}},{"description":"linter-b: another issue","severity":"error","fingerprint":"0777B4FE60242BD8B2E9B7E92C4B9521","location":{"path":"path/to/fileb.go","lines":{"begin":300}}},{"description":"linter-c: issue c","severity":"critical","fingerprint":"BEE6E9FBB6BFA4B7DB9FB036697FB036","location":{"path":"path/to/filec.go","lines":{"begin":200}}}]
`

	assert.Equal(t, expected, buf.String())
}
