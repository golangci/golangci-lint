package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestSarif_Print(t *testing.T) {
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
			FromLinter: "linter-a",
			Severity:   "error",
			Text:       "some issue 2",
			Pos: token.Position{
				Filename: "path/to/filec.go",
				Offset:   3,
				Line:     11,
				Column:   5,
			},
		},
	}

	buf := new(bytes.Buffer)

	printer := NewSarif(nil, buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `{"version":"2.1.0","$schema":"https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.4.json","runs":[{"tool":{"driver":{"name":"linter-a"}},"results":[{"ruleId":"some issue","level":"warning","message":{"text":"some issue"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/filea.go","index":0},"region":{"startLine":10,"startColumn":4}}}]},{"ruleId":"some issue 2","level":"error","message":{"text":"some issue 2"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/filec.go","index":0},"region":{"startLine":11,"startColumn":5}}}]}]},{"tool":{"driver":{"name":"linter-b"}},"results":[{"ruleId":"another issue","level":"error","message":{"text":"another issue"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/fileb.go","index":0},"region":{"startLine":300,"startColumn":9}}}]}]}]}
`

	assert.Equal(t, expected, buf.String())
}
