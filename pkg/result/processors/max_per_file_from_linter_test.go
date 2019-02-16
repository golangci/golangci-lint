package processors

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

func newFromLinterIssue(linterName string) result.Issue {
	return result.Issue{
		FromLinter: linterName,
	}
}

func TestMaxPerFileFromLinterUnlimited(t *testing.T) {
	p := NewMaxPerFileFromLinter(&config.Config{})
	gosimple := newFromLinterIssue("gosimple")
	processAssertSame(t, p, gosimple) // collect stat
	processAssertSame(t, p, gosimple) // check not limits
}

func TestMaxPerFileFromLinter(t *testing.T) {
	p := NewMaxPerFileFromLinter(&config.Config{})
	for _, name := range []string{"gofmt", "goimports"} {
		limited := newFromLinterIssue(name)
		gosimple := newFromLinterIssue("gosimple")
		processAssertSame(t, p, limited)
		processAssertSame(t, p, gosimple)
		processAssertEmpty(t, p, limited)
	}
}
