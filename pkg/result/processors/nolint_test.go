package processors

import (
	"go/token"
	"path/filepath"
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
)

func newNolintFileIssue(line int, fromLinter string) result.Issue {
	return result.Issue{
		File:       filepath.Join("testdata", "nolint.go"),
		LineNumber: line,
		FromLinter: fromLinter,
	}
}

func TestNolint(t *testing.T) {
	p := NewNolint(token.NewFileSet())
	processAssertEmpty(t, p, newNolintFileIssue(3, "gofmt"))
	processAssertEmpty(t, p, newNolintFileIssue(3, "gofmt")) // check cached is ok
	processAssertSame(t, p, newNolintFileIssue(3, "gofmtA")) // check different name

	processAssertEmpty(t, p, newNolintFileIssue(4, "any"))
	processAssertEmpty(t, p, newNolintFileIssue(5, "any"))

	processAssertSame(t, p, newNolintFileIssue(1, "golint"))
}
