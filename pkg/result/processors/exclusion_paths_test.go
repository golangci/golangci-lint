package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestExclusionPaths_Process(t *testing.T) {
	logger := logutils.NewStderrLog(logutils.DebugKeyEmpty)
	logger.SetLevel(logutils.LogLevelDebug)

	testCases := []struct {
		desc     string
		patterns []string
		issues   []result.Issue
		expected []result.Issue
	}{
		{
			desc:     "word",
			patterns: []string{"foo"},
			issues: []result.Issue{
				{RelativePath: "foo.go"},
				{RelativePath: "foo/foo.go"},
				{RelativePath: "foo/bar.go"},
				{RelativePath: "bar/foo.go"},
				{RelativePath: "bar/bar.go"},
			},
			expected: []result.Issue{
				{RelativePath: "bar/bar.go"},
			},
		},
		{
			desc:     "begin with word",
			patterns: []string{"^foo"},
			issues: []result.Issue{
				{RelativePath: "foo.go"},
				{RelativePath: "foo/foo.go"},
				{RelativePath: "foo/bar.go"},
				{RelativePath: "bar/foo.go"},
				{RelativePath: "bar/bar.go"},
			},
			expected: []result.Issue{
				{RelativePath: "bar/foo.go"},
				{RelativePath: "bar/bar.go"},
			},
		},
		{
			desc:     "directory begin with word",
			patterns: []string{"^foo/"},
			issues: []result.Issue{
				{RelativePath: "foo.go"},
				{RelativePath: "foo/foo.go"},
				{RelativePath: "foo/bar.go"},
				{RelativePath: "bar/foo.go"},
				{RelativePath: "bar/bar.go"},
			},
			expected: []result.Issue{
				{RelativePath: "foo.go"},
				{RelativePath: "bar/foo.go"},
				{RelativePath: "bar/bar.go"},
			},
		},
		{
			desc:     "same suffix with unconstrained expression",
			patterns: []string{"c/d.go"},
			issues: []result.Issue{
				{RelativePath: "a/b/c/d.go"},
				{RelativePath: "c/d.go"},
			},
			expected: []result.Issue{},
		},
		{
			desc:     "same suffix with constrained expression",
			patterns: []string{"^c/d.go"},
			issues: []result.Issue{
				{RelativePath: "a/b/c/d.go"},
				{RelativePath: "c/d.go"},
			},
			expected: []result.Issue{
				{RelativePath: "a/b/c/d.go"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p, err := NewExclusionPaths(logger, &config.LinterExclusions{Paths: test.patterns})
			require.NoError(t, err)

			processedIssues := process(t, p, test.issues...)

			assert.Equal(t, test.expected, processedIssues)

			p.Finish()
		})
	}
}
