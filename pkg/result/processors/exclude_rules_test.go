package processors

import (
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestExcludeRulesMultiple(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
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
		{
			Source:  "^//go:generate ",
			Linters: []string{"lll"},
		},
	}, lineCache, nil)

	cases := []issueCase{
		{Path: "e.go", Text: "exclude", Linter: "linter"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "testonly", Linter: "linter"},
		{Path: filepath.Join("testdata", "exclude_rules.go"), Line: 3, Linter: "lll"},
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
			Line:   i.Line(),
		})
	}
	expectedCases := []issueCase{
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
	}
	assert.Equal(t, expectedCases, resultingCases)
}

type issueCase struct {
	Path   string
	Line   int
	Text   string
	Linter string
}

func newIssueCase(c issueCase) result.Issue {
	return result.Issue{
		Text:       c.Text,
		FromLinter: c.Linter,
		Pos: token.Position{
			Filename: c.Path,
			Line:     c.Line,
		},
	}
}

func TestExcludeRulesText(t *testing.T) {
	p := NewExcludeRules([]ExcludeRule{
		{
			Text: "^exclude$",
			Linters: []string{
				"linter",
			},
		},
	}, nil, nil)
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
}

func TestExcludeRulesEmpty(t *testing.T) {
	processAssertSame(t, NewExcludeRules(nil, nil, nil), newTextIssue("test"))
}

func TestExcludeRulesCaseSensitiveMultiple(t *testing.T) {
	lineCache := fsutils.NewLineCache(fsutils.NewFileCache())
	p := NewExcludeRulesCaseSensitive([]ExcludeRule{
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
		{
			Source:  "^//go:generate ",
			Linters: []string{"lll"},
		},
	}, lineCache, nil)

	cases := []issueCase{
		{Path: "e.go", Text: "exclude", Linter: "linter"},
		{Path: "e.go", Text: "excLude", Linter: "linter"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "testonly", Linter: "linter"},
		{Path: "e_test.go", Text: "testOnly", Linter: "linter"},
		{Path: filepath.Join("testdata", "exclude_rules_case_sensitive.go"), Line: 3, Linter: "lll"},
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
			Line:   i.Line(),
		})
	}
	expectedCases := []issueCase{
		{Path: "e.go", Text: "excLude", Linter: "linter"},
		{Path: "e.go", Text: "some", Linter: "linter"},
		{Path: "e_Test.go", Text: "normal", Linter: "testlinter"},
		{Path: "e_test.go", Text: "another", Linter: "linter"},
		{Path: "e_test.go", Text: "testOnly", Linter: "linter"},
		{Path: filepath.Join("testdata", "exclude_rules_case_sensitive.go"), Line: 3, Linter: "lll"},
	}
	assert.Equal(t, expectedCases, resultingCases)
}

func TestExcludeRulesCaseSensitiveText(t *testing.T) {
	p := NewExcludeRulesCaseSensitive([]ExcludeRule{
		{
			Text: "^exclude$",
			Linters: []string{
				"linter",
			},
		},
	}, nil, nil)
	texts := []string{"exclude", "excLude", "1", "", "exclud", "notexclude"}
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
}

func TestExcludeRulesCaseSensitiveEmpty(t *testing.T) {
	processAssertSame(t, NewExcludeRulesCaseSensitive(nil, nil, nil), newTextIssue("test"))
}
