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
			Severity:   "low",
			Text:       "some issue 2",
			Pos: token.Position{
				Filename: "path/to/filec.go",
				Offset:   3,
				Line:     11,
				Column:   5,
			},
		},
		{
			FromLinter: "linter-c",
			Severity:   "error",
			Text:       "some issue without column",
			Pos: token.Position{
				Filename: "path/to/filed.go",
				Offset:   3,
				Line:     11,
			},
		},
	}

	buf := new(bytes.Buffer)

	printer := NewSarif(buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `{"version":"2.1.0","$schema":"https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.6.json","runs":[{"tool":{"driver":{"name":"golangci-lint"}},"results":[{"ruleId":"linter-a","level":"warning","message":{"text":"some issue"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/filea.go","index":0},"region":{"startLine":10,"startColumn":4}}}]},{"ruleId":"linter-b","level":"error","message":{"text":"another issue"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/fileb.go","index":0},"region":{"startLine":300,"startColumn":9}}}]},{"ruleId":"linter-a","level":"error","message":{"text":"some issue 2"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/filec.go","index":0},"region":{"startLine":11,"startColumn":5}}}]},{"ruleId":"linter-c","level":"error","message":{"text":"some issue without column"},"locations":[{"physicalLocation":{"artifactLocation":{"uri":"path/to/filed.go","index":0},"region":{"startLine":11,"startColumn":1}}}]}]}]}
`

	assert.Equal(t, expected, buf.String())
}

func TestSarif_Print_empty(t *testing.T) {
	buf := new(bytes.Buffer)

	printer := NewSarif(buf)

	err := printer.Print(nil)
	require.NoError(t, err)

	expected := `{"version":"2.1.0","$schema":"https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.6.json","runs":[{"tool":{"driver":{"name":"golangci-lint"}},"results":[]}]}
`

	assert.Equal(t, expected, buf.String())
}
