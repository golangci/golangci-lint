package processors

import (
	"go/token"
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
)

func newULIssue(file string, line int) result.Issue {
	return result.Issue{
		Pos: token.Position{
			Filename: file,
			Line:     line,
		},
	}
}

func TestUniqByLine(t *testing.T) {
	p := NewUniqByLine(true)
	i1 := newULIssue("f1", 1)

	processAssertSame(t, p, i1)
	processAssertEmpty(t, p, i1) // check skipping
	processAssertEmpty(t, p, i1) // check accumulated error

	processAssertSame(t, p, newULIssue("f1", 2)) // another line
	processAssertSame(t, p, newULIssue("f2", 1)) // another file
}

func TestUniqByLineDisabled(t *testing.T) {
	p := NewUniqByLine(false)
	i1 := newULIssue("f1", 1)

	processAssertSame(t, p, i1)
	processAssertSame(t, p, i1) // check the same issue passed twice
}
