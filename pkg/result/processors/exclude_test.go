package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestExclude(t *testing.T) {
	p := NewExclude(&config.Issues{ExcludePatterns: []string{"^exclude$"}})

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

func TestExclude_empty(t *testing.T) {
	processAssertSame(t, NewExclude(&config.Issues{}), newIssueFromTextTestCase("test"))
}

func TestExclude_caseSensitive(t *testing.T) {
	p := NewExclude(&config.Issues{ExcludePatterns: []string{"^exclude$"}, ExcludeCaseSensitive: true})

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
