package processors

import (
	"testing"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func TestMaxFromLinter(t *testing.T) {
	p := NewMaxFromLinter(1, logutils.NewStderrLog(logutils.DebugKeyEmpty), &config.Config{})
	gosimple := newFromLinterIssue("gosimple")
	gofmt := newFromLinterIssue("gofmt")
	processAssertSame(t, p, gosimple)  // ok
	processAssertSame(t, p, gofmt)     // ok: another
	processAssertEmpty(t, p, gosimple) // skip
}
