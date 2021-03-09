package processors

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestPathPrefixer_Process(t *testing.T) {
	paths := func(ps ...string) (issues []result.Issue) {
		for _, p := range ps {
			issues = append(issues, result.Issue{Pos: token.Position{Filename: p}})
		}
		return
	}
	for _, tt := range []struct {
		name, prefix string
		issues, want []result.Issue
	}{
		{"empty prefix", "", paths("some/path", "cool"), paths("some/path", "cool")},
		{"prefix", "ok", paths("some/path", "cool"), paths("ok/some/path", "ok/cool")},
		{"prefix slashed", "ok/", paths("some/path", "cool"), paths("ok/some/path", "ok/cool")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p := NewPathPrefixer(tt.prefix)
			got, err := p.Process(tt.issues)
			r.NoError(err, "prefixer should never error")

			r.Equal(got, tt.want)
		})
	}
}
