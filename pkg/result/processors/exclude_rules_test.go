package processors

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestExcludeRules(t *testing.T) {
	t.Run("Multiple", func(t *testing.T) {
		p := NewExcludeRules([]ExcludeRule{
			{
				Text:    "^exclude$",
				Linters: []string{"linter"},
			},
			{
				Linters: []string{"testlinter"},
				Path:    `_test\.go`,
			},
			{
				Text: "^testonly$",
				Path: `_test\.go`,
			},
		})
		type issueCase struct {
			Path   string
			Text   string
			Linter string
		}
		var newIssueCase = func(c issueCase) result.Issue {
			return result.Issue{
				Text:       c.Text,
				FromLinter: c.Linter,
				Pos: token.Position{
					Filename: c.Path,
				},
			}
		}
		cases := []issueCase{
			{Path: "e.go", Text: "exclude", Linter: "linter"},
			{Path: "e.go", Text: "some", Linter: "linter"},
			{Path: "e_test.go", Text: "normal", Linter: "testlinter"},
			{Path: "e_test.go", Text: "another", Linter: "linter"},
			{Path: "e_test.go", Text: "testonly", Linter: "linter"},
		}
		var issues []result.Issue
		for _, c := range cases {
			issues = append(issues, newIssueCase(c))
		}
		processedIssues := process(t, p, issues...)
		var resultingCases []issueCase
		for _, i := range processedIssues {
			resultingCases = append(resultingCases, issueCase{
				Path:   i.FilePath(),
				Linter: i.FromLinter,
				Text:   i.Text,
			})
		}
		expectedCases := []issueCase{
			{Path: "e.go", Text: "some", Linter: "linter"},
			{Path: "e_test.go", Text: "another", Linter: "linter"},
		}
		assert.Equal(t, expectedCases, resultingCases)
	})
	t.Run("Text", func(t *testing.T) {
		p := NewExcludeRules([]ExcludeRule{
			{
				Text: "^exclude$",
				Linters: []string{
					"linter",
				},
			},
		})
		texts := []string{"excLude", "1", "", "exclud", "notexclude"}
		var issues []result.Issue
		for _, t := range texts {
			issues = append(issues, result.Issue{
				Text:       t,
				FromLinter: "linter",
			})
		}

		processedIssues := process(t, p, issues...)
		assert.Len(t, processedIssues, len(issues)-1)

		var processedTexts []string
		for _, i := range processedIssues {
			processedTexts = append(processedTexts, i.Text)
		}
		assert.Equal(t, texts[1:], processedTexts)
	})
	t.Run("Empty", func(t *testing.T) {
		processAssertSame(t, NewExcludeRules(nil), newTextIssue("test"))
	})
}
