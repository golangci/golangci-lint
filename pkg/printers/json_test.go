package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestJSON_Print(t *testing.T) {
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

	printer := NewJSON(nil, buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `{"Issues":[{"FromLinter":"linter-a","Text":"some issue","Severity":"warning","SourceLines":null,"Replacement":null,"Pos":{"Filename":"path/to/filea.go","Offset":2,"Line":10,"Column":4},"ExpectNoLint":false,"ExpectedNoLintLinter":""},{"FromLinter":"linter-b","Text":"another issue","Severity":"error","SourceLines":["func foo() {","\tfmt.Println(\"bar\")","}"],"Replacement":null,"Pos":{"Filename":"path/to/fileb.go","Offset":5,"Line":300,"Column":9},"ExpectNoLint":false,"ExpectedNoLintLinter":""}],"Report":null}
`

	assert.Equal(t, expected, buf.String())
}
