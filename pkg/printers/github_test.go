//nolint:dupl
package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestGithub_Print(t *testing.T) {
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
	printer := NewGithub(buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `::warning file=path/to/filea.go,line=10,col=4::some issue (linter-a)
::error file=path/to/fileb.go,line=300,col=9::another issue (linter-b)
`

	assert.Equal(t, expected, buf.String())
}

func TestFormatGithubIssue(t *testing.T) {
	sampleIssue := result.Issue{
		FromLinter: "sample-linter",
		Text:       "some issue",
		Pos: token.Position{
			Filename: "path/to/file.go",
			Offset:   2,
			Line:     10,
			Column:   4,
		},
	}
	require.Equal(t, "::error file=path/to/file.go,line=10,col=4::some issue (sample-linter)", formatIssueAsGithub(&sampleIssue))

	sampleIssue.Pos.Column = 0
	require.Equal(t, "::error file=path/to/file.go,line=10::some issue (sample-linter)", formatIssueAsGithub(&sampleIssue))
}
