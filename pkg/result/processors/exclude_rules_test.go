package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestExcludeRules(t *testing.T) {
	p := NewExcludeRules([]ExcludeRule{
		{
			Text: "^exclude$",
		},
	})
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
	t.Run("Empty", func(t *testing.T) {
		processAssertSame(t, NewExcludeRules(nil), newTextIssue("test"))
	})
}
