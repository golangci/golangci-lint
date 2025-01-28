package processors

import (
	"path/filepath"
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
		cfg      *config.LinterExclusions
		issues   []result.Issue
		expected []result.Issue
	}{
		{
			desc: "paths: word",
			cfg: &config.LinterExclusions{
				Paths: []string{"foo"},
			},
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
			desc: "paths: begin with word",
			cfg: &config.LinterExclusions{
				Paths: []string{"^foo"},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("foo.go")},
				{RelativePath: filepath.FromSlash("foo/foo.go")},
				{RelativePath: filepath.FromSlash("foo/bar.go")},
				{RelativePath: filepath.FromSlash("bar/foo.go")},
				{RelativePath: filepath.FromSlash("bar/bar.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("bar/foo.go")},
				{RelativePath: filepath.FromSlash("bar/bar.go")},
			},
		},
		{
			desc: "paths: directory begin with word",
			cfg: &config.LinterExclusions{
				Paths: []string{"^foo/"},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("foo.go")},
				{RelativePath: filepath.FromSlash("foo/foo.go")},
				{RelativePath: filepath.FromSlash("foo/bar.go")},
				{RelativePath: filepath.FromSlash("bar/foo.go")},
				{RelativePath: filepath.FromSlash("bar/bar.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("foo.go")},
				{RelativePath: filepath.FromSlash("bar/foo.go")},
				{RelativePath: filepath.FromSlash("bar/bar.go")},
			},
		},
		{
			desc: "paths: same suffix with unconstrained expression",
			cfg: &config.LinterExclusions{
				Paths: []string{"c/d.go"},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("a/b/c/d.go")},
				{RelativePath: filepath.FromSlash("c/d.go")},
			},
			expected: []result.Issue{},
		},
		{
			desc: "paths: same suffix with constrained expression",
			cfg: &config.LinterExclusions{
				Paths: []string{"^c/d.go"},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("a/b/c/d.go")},
				{RelativePath: filepath.FromSlash("c/d.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("a/b/c/d.go")},
			},
		},
		{
			desc: "paths: unused",
			cfg: &config.LinterExclusions{
				WarnUnused: true,
				Paths: []string{
					`^z/d.go`, // This pattern is unused.
					`^c/d.go`,
				},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("a/b/c/d.go")},
				{RelativePath: filepath.FromSlash("c/d.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("a/b/c/d.go")},
			},
		},
		{
			desc: "pathsExcept",
			cfg: &config.LinterExclusions{
				PathsExcept: []string{`^base/c/.*$`},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/c/file.go")},
				{RelativePath: filepath.FromSlash("base/c/a/file.go")},
				{RelativePath: filepath.FromSlash("base/c/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
		},
		{
			desc: "pathsExcept: unused",
			cfg: &config.LinterExclusions{
				WarnUnused: true,
				PathsExcept: []string{
					`^base/z/.*$`, // This pattern is unused.
					`^base/c/.*$`,
				},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/c/file.go")},
				{RelativePath: filepath.FromSlash("base/c/a/file.go")},
				{RelativePath: filepath.FromSlash("base/c/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
		},
		{
			desc: "pathsExcept: multiple patterns",
			cfg: &config.LinterExclusions{
				PathsExcept: []string{
					`^base/e/.*$`,
					`^base/c/.*$`,
				},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/c/file.go")},
				{RelativePath: filepath.FromSlash("base/c/a/file.go")},
				{RelativePath: filepath.FromSlash("base/c/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
				{RelativePath: filepath.FromSlash("base/e/file.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
		},
		{
			desc: "pathsExcept and paths",
			cfg: &config.LinterExclusions{
				Paths:       []string{"^base/b/"},
				PathsExcept: []string{`^base/c/.*$`},
			},
			issues: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/b/file.go")},
				{RelativePath: filepath.FromSlash("base/c/file.go")},
				{RelativePath: filepath.FromSlash("base/c/a/file.go")},
				{RelativePath: filepath.FromSlash("base/c/b/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
			expected: []result.Issue{
				{RelativePath: filepath.FromSlash("base/a/file.go")},
				{RelativePath: filepath.FromSlash("base/d/file.go")},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p, err := NewExclusionPaths(logger, test.cfg)
			require.NoError(t, err)

			processedIssues := process(t, p, test.issues...)

			assert.Equal(t, test.expected, processedIssues)

			p.Finish()
		})
	}
}
