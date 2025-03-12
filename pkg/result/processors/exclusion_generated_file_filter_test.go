package processors

import (
	"fmt"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

func TestGeneratedFileFilter_shouldPassIssue(t *testing.T) {
	testCases := []struct {
		desc   string
		mode   string
		issue  *result.Issue
		assert assert.BoolAssertionFunc
	}{
		{
			desc: "lax ",
			mode: config.GeneratedModeLax,
			issue: &result.Issue{
				FromLinter: "example",
				Pos: token.Position{
					Filename: filepath.FromSlash("testdata/exclusion_generated_file_filter/go_strict_invalid.go"),
				},
			},
			assert: assert.False,
		},
		{
			desc: "strict ",
			mode: config.GeneratedModeStrict,
			issue: &result.Issue{
				FromLinter: "example",
				Pos: token.Position{
					Filename: filepath.FromSlash("testdata/exclusion_generated_file_filter/go_strict_invalid.go"),
				},
			},
			assert: assert.True,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p := NewGeneratedFileFilter(test.mode)

			pass, err := p.shouldPassIssue(test.issue)
			require.NoError(t, err)

			test.assert(t, pass)
		})
	}
}

func TestGeneratedFileFilter_shouldPassIssue_error(t *testing.T) {
	notFoundMsg := "no such file or directory"
	if runtime.GOOS == "windows" {
		notFoundMsg = "The system cannot find the file specified."
	}

	testCases := []struct {
		desc     string
		mode     string
		issue    *result.Issue
		expected string
	}{
		{
			desc: "non-existing file (lax)",
			mode: config.GeneratedModeLax,
			issue: &result.Issue{
				FromLinter: "example",
				Pos: token.Position{
					Filename: filepath.FromSlash("no-existing.go"),
				},
			},
			expected: fmt.Sprintf("failed to get doc (lax) of file %[1]s: failed to parse file: open %[1]s: %[2]s",
				filepath.FromSlash("no-existing.go"), notFoundMsg),
		},
		{
			desc: "non-existing file (strict)",
			mode: config.GeneratedModeStrict,
			issue: &result.Issue{
				FromLinter: "example",
				Pos: token.Position{
					Filename: filepath.FromSlash("no-existing.go"),
				},
			},
			expected: fmt.Sprintf("failed to get doc (strict) of file %[1]s: failed to parse file: open %[1]s: %[2]s",
				filepath.FromSlash("no-existing.go"), notFoundMsg),
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p := NewGeneratedFileFilter(test.mode)

			pass, err := p.shouldPassIssue(test.issue)

			//nolint:testifylint // It's a loop and the main expectation is the error message.
			assert.EqualError(t, err, test.expected)
			assert.False(t, pass)
		})
	}
}
