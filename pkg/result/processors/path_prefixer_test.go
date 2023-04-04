package processors

import (
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestPathPrefixer_Process(t *testing.T) {
	paths := func(ps ...string) (issues []result.Issue) {
		for _, p := range ps {
			issues = append(issues, result.Issue{Pos: token.Position{Filename: filepath.FromSlash(p)}})
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
			want:   paths("some/path", "cool"),
		},
		{
			name:   "prefix",
			prefix: "ok",
			issues: paths("some/path", "cool"),
			want:   paths("ok/some/path", "ok/cool"),
		},
		{
			name:   "prefix slashed",
			prefix: "ok/",
			issues: paths("some/path", "cool"),
			want:   paths("ok/some/path", "ok/cool"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPathPrefixer(tt.prefix)

			got, err := p.Process(tt.issues)
			require.NoError(t, err)

			assert.Equal(t, got, tt.want)
		})
	}
}
