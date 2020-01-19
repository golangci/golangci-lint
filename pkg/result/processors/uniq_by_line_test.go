package processors

import (
	"go/token"
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

func newFLIssue(file string, line int) result.Issue {
	return result.Issue{
		Pos: token.Position{
			Filename: file,
			Line:     line,
		},
	}
}

func TestUniqByLine(t *testing.T) {
	cfg := config.Config{}
	cfg.Output.UniqByLine = true

	p := NewUniqByLine(&cfg)
	i1 := newFLIssue("f1", 1)

	processAssertSame(t, p, i1)
	processAssertEmpty(t, p, i1) // check skipping
	processAssertEmpty(t, p, i1) // check accumulated error

	processAssertSame(t, p, newFLIssue("f1", 2)) // another line
	processAssertSame(t, p, newFLIssue("f2", 1)) // another file
}

func TestUniqByLineDisabled(t *testing.T) {
	cfg := config.Config{}
	cfg.Output.UniqByLine = false

	p := NewUniqByLine(&cfg)
	i1 := newFLIssue("f1", 1)

	processAssertSame(t, p, i1)
	processAssertSame(t, p, i1) // check the same issue passed twice
}
