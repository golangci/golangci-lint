package processors

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/stretchr/testify/assert"
)

func newTextIssue(text string) result.Issue {
	return result.Issue{
		Text: text,
	}
}

func process(t *testing.T, p Processor, issues ...result.Issue) []result.Issue {
	processedIssues, err := p.Process(issues)
	assert.NoError(t, err)
	return processedIssues
}

func processAssertSame(t *testing.T, p Processor, issues ...result.Issue) {
	processedIssues := process(t, p, issues...)
	assert.Equal(t, issues, processedIssues)
}

func processAssertEmpty(t *testing.T, p Processor, issues ...result.Issue) {
	processedIssues := process(t, p, issues...)
	assert.Empty(t, processedIssues)
}

func TestExclude(t *testing.T) {
	p := NewExclude("^exclude$")
	texts := []string{"excLude", "1", "", "exclud", "notexclude"}
	var issues []result.Issue
	for _, t := range texts {
		issues = append(issues, newTextIssue(t))
	}

	processedIssues := process(t, p, issues...)
	assert.Len(t, processedIssues, len(issues)-1)

	var processedTexts []string
	for _, i := range processedIssues {
		processedTexts = append(processedTexts, i.Text)
	}
	assert.Equal(t, texts[1:], processedTexts)
}

func TestNoExclude(t *testing.T) {
	processAssertSame(t, NewExclude(""), newTextIssue("test"))
}
