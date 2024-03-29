package processors

import (
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func newFileIssue(file string) result.Issue {
	return result.Issue{
		Pos: token.Position{
			Filename: file,
		},
	}
}

func newTestSkipFiles(t *testing.T, patterns ...string) *SkipFiles {
	p, err := NewSkipFiles(patterns, "")
	require.NoError(t, err)
	return p
}

func TestSkipFiles(t *testing.T) {
	processAssertSame(t, newTestSkipFiles(t), newFileIssue("any.go"))

	processAssertEmpty(t, newTestSkipFiles(t, "file"),
		newFileIssue("file.go"),
		newFileIssue("file"),
		newFileIssue("nofile.go"))

	processAssertEmpty(t, newTestSkipFiles(t, ".*"), newFileIssue("any.go"))

	cleanPath := strings.ReplaceAll(filepath.FromSlash("a/b/c.go"), `\`, `\\`)
	processAssertEmpty(t, newTestSkipFiles(t, cleanPath), newFileIssue(filepath.FromSlash("a/b/c.go")))
	processAssertSame(t, newTestSkipFiles(t, cleanPath), newFileIssue(filepath.FromSlash("a/b/d.go")))

	processAssertEmpty(t, newTestSkipFiles(t, ".*\\.pb\\.go"), newFileIssue(filepath.FromSlash("a/b.pb.go")))
	processAssertSame(t, newTestSkipFiles(t, ".*\\.pb\\.go"), newFileIssue(filepath.FromSlash("a/b.go")))

	processAssertEmpty(t, newTestSkipFiles(t, ".*\\.pb\\.go$"), newFileIssue(filepath.FromSlash("a/b.pb.go")))
	processAssertSame(t, newTestSkipFiles(t, ".*\\.pb\\.go$"), newFileIssue(filepath.FromSlash("a/b.go")))
}

func TestSkipFilesInvalidPattern(t *testing.T) {
	p, err := NewSkipFiles([]string{"\\o"}, "")
	require.Error(t, err)
	assert.Nil(t, p)
}
