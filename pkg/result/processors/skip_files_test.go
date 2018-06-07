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
	p := newTestSkipFiles(t)
	processAssertSame(t, p, newFileIssue("any.go"))

	p = newTestSkipFiles(t, "file")
	processAssertEmpty(t, p,
		newFileIssue("file.go"),
		newFileIssue("file"),
		newFileIssue("nofile.go"))

	p = newTestSkipFiles(t, ".*")
	processAssertEmpty(t, p, newFileIssue("any.go"))
}

func TestSkipFilesInvalidPattern(t *testing.T) {
	p, err := NewSkipFiles([]string{"\\o"})
	assert.Error(t, err)
	assert.Nil(t, p)
}
