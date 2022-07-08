package processors

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/stretchr/testify/assert"
)

func TestExclude(t *testing.T) {
	p := NewExclude("^exclude$")
	texts := []string{"excLude", "1", "", "exclud", "notexclude"}
	var issues []result.Issue
	for _, t := range texts {
		issues = append(issues, newIssueFromTextTestCase(t))
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
	processAssertSame(t, NewExclude(""), newIssueFromTextTestCase("test"))
}

func TestExcludeCaseSensitive(t *testing.T) {
	p := NewExcludeCaseSensitive("^exclude$")
	texts := []string{"excLude", "1", "", "exclud", "exclude"}
	var issues []result.Issue
	for _, t := range texts {
		issues = append(issues, newIssueFromTextTestCase(t))
	}

	processedIssues := process(t, p, issues...)
	assert.Len(t, processedIssues, len(issues)-1)

	var processedTexts []string
	for _, i := range processedIssues {
		processedTexts = append(processedTexts, i.Text)
	}
	assert.Equal(t, texts[:len(texts)-1], processedTexts)
}
