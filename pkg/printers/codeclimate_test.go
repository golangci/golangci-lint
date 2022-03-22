//nolint:dupl
package printers

import (
	"bytes"
	"context"
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
	}

	buf := new(bytes.Buffer)
	printer := NewCodeClimate(buf)

	err := printer.Print(context.Background(), issues)
	require.NoError(t, err)

	//nolint:lll
	expected := `[{"description":"linter-a: some issue","severity":"warning","fingerprint":"BA73C5DF4A6FD8462FFF1D3140235777","location":{"path":"path/to/filea.go","lines":{"begin":10}}},{"description":"linter-b: another issue","severity":"error","fingerprint":"0777B4FE60242BD8B2E9B7E92C4B9521","location":{"path":"path/to/fileb.go","lines":{"begin":300}}}]`

	assert.Equal(t, expected, buf.String())
}
