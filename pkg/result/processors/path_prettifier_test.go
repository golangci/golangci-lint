package processors

import (
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func newPPIssue(fn, rp string) result.Issue {
	return result.Issue{
		Pos:          token.Position{Filename: filepath.FromSlash(fn)},
		RelativePath: filepath.FromSlash(rp),
	}
}

func TestPathPrettifier_Process(t *testing.T) {
	paths := func(ps ...string) (issues []result.Issue) {
		for _, p := range ps {
			issues = append(issues, newPPIssue("test", p))
		}
		return
	}

	for _, tt := range []struct {
		name, prefix string
		issues, want []result.Issue
	}{
		{
			name:   "empty prefix",
			issues: paths("some/path", "cool"),
			want: []result.Issue{
				newPPIssue("some/path", "some/path"),
				newPPIssue("cool", "cool"),
			},
		},
		{
			name:   "prefix",
			prefix: "ok",
			issues: paths("some/path", "cool"),
			want: []result.Issue{
				newPPIssue("ok/some/path", "some/path"),
				newPPIssue("ok/cool", "cool"),
			},
		},
		{
			name:   "prefix slashed",
			prefix: "ok/",
			issues: paths("some/path", "cool"),
			want: []result.Issue{
				newPPIssue("ok/some/path", "some/path"),
				newPPIssue("ok/cool", "cool"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPathPrettifier(logutils.NewStderrLog(logutils.DebugKeyEmpty), tt.prefix)

			got, err := p.Process(tt.issues)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
