package processors

import (
	"go/token"
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/stretchr/testify/assert"
)

func newFileIssue(file string) result.Issue {
	return result.Issue{
		Pos: token.Position{
			Filename: file,
		},
	}
}

func newTestSkipFiles(t *testing.T, patterns ...string) *SkipFiles {
	p, err := NewSkipFiles(patterns)
	assert.NoError(t, err)
	return p
}

func TestSkipFiles(t *testing.T) {
	processAssertSame(t, newTestSkipFiles(t), newFileIssue("any.go"))

	processAssertEmpty(t, newTestSkipFiles(t, "file"),
		newFileIssue("file.go"),
		newFileIssue("file"),
		newFileIssue("nofile.go"))

	processAssertEmpty(t, newTestSkipFiles(t, ".*"), newFileIssue("any.go"))

	processAssertEmpty(t, newTestSkipFiles(t, "a/b/c.go"), newFileIssue("a/b/c.go"))
	processAssertSame(t, newTestSkipFiles(t, "a/b/c.go"), newFileIssue("a/b/d.go"))

	processAssertEmpty(t, newTestSkipFiles(t, ".*\\.pb\\.go"), newFileIssue("a/b.pb.go"))
	processAssertSame(t, newTestSkipFiles(t, ".*\\.pb\\.go"), newFileIssue("a/b.go"))

	processAssertEmpty(t, newTestSkipFiles(t, ".*\\.pb\\.go$"), newFileIssue("a/b.pb.go"))
	processAssertSame(t, newTestSkipFiles(t, ".*\\.pb\\.go$"), newFileIssue("a/b.go"))
}

func TestSkipFilesInvalidPattern(t *testing.T) {
	p, err := NewSkipFiles([]string{"\\o"})
	assert.Error(t, err)
	assert.Nil(t, p)
}
