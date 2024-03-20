package processors

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestInvalidIssue_Process(t *testing.T) {
	logger := logutils.NewStderrLog(logutils.DebugKeyInvalidIssue)
	logger.SetLevel(logutils.LogLevelDebug)

	p := NewInvalidIssue(logger)

	testCases := []struct {
		desc     string
		issues   []result.Issue
		expected []result.Issue
	}{
		{
			desc: "typecheck",
			issues: []result.Issue{
				{FromLinter: "typecheck"},
			},
			expected: []result.Issue{
				{FromLinter: "typecheck"},
			},
		},
		{
			desc: "Go file",
			issues: []result.Issue{
				{
					FromLinter: "example",
					Pos: token.Position{
						Filename: "test.go",
					},
				},
			},
			expected: []result.Issue{
				{
					FromLinter: "example",
					Pos: token.Position{
						Filename: "test.go",
					},
				},
			},
		},
		{
			desc: "go.mod",
			issues: []result.Issue{
				{
					FromLinter: "example",
					Pos: token.Position{
						Filename: "go.mod",
					},
				},
			},
			expected: []result.Issue{
				{
					FromLinter: "example",
					Pos: token.Position{
						Filename: "go.mod",
					},
				},
			},
		},
		{
			desc: "non Go file",
			issues: []result.Issue{
				{
					FromLinter: "example",
					Pos: token.Position{
						Filename: "test.txt",
					},
				},
			},
			expected: []result.Issue{},
		},
		{
			desc: "no filename",
			issues: []result.Issue{
				{
					FromLinter: "example",
					Pos: token.Position{
						Filename: "",
					},
				},
			},
			expected: []result.Issue{},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			after, err := p.Process(test.issues)
			require.NoError(t, err)

			assert.Equal(t, test.expected, after)
		})
	}
}
