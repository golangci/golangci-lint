package printers

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

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
