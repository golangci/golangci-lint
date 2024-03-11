package processors

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

type issueTestCase struct {
	Path     string
	Line     int
	Text     string
	Linter   string
	Severity string
}

func newIssueFromIssueTestCase(c issueTestCase) result.Issue {
	return result.Issue{
		Text:       c.Text,
		FromLinter: c.Linter,
		Severity:   c.Severity,
		Pos: token.Position{
			Filename: c.Path,
			Line:     c.Line,
		},
	}
}

func newIssueFromTextTestCase(text string) result.Issue {
	return result.Issue{
		Text: text,
	}
}

func process(t *testing.T, p Processor, issues ...result.Issue) []result.Issue {
	t.Helper()

	processedIssues, err := p.Process(issues)
	require.NoError(t, err)
	return processedIssues
}

func processAssertSame(t *testing.T, p Processor, issues ...result.Issue) {
	t.Helper()

	processedIssues := process(t, p, issues...)
	assert.Equal(t, issues, processedIssues)
}

func processAssertEmpty(t *testing.T, p Processor, issues ...result.Issue) {
	t.Helper()

	processedIssues := process(t, p, issues...)
	assert.Empty(t, processedIssues)
}
